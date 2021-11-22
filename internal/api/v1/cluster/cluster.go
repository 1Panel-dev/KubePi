package cluster

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterrepo"
	"strings"
	"sync"
	"time"

	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/KubeOperator/kubepi/pkg/certificate"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	authV1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handler struct {
	clusterService        cluster.Service
	clusterBindingService clusterbinding.Service
	clusterRepoService    clusterrepo.Service
}

func NewHandler() *Handler {
	return &Handler{
		clusterService:        cluster.NewService(),
		clusterBindingService: clusterbinding.NewService(),
		clusterRepoService:    clusterrepo.NewService(),
	}
}

func (h *Handler) CreateCluster() iris.Handler {
	return func(ctx *context.Context) {
		var req Cluster
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if req.ConfigFileContentStr != "" {
			req.Spec.Authentication.ConfigFileContent = []byte(req.ConfigFileContentStr)
		}
		if req.CaDataStr != "" {
			req.CaCertificate.CertData = []byte(req.CaDataStr)
		}
		if req.Spec.Authentication.Mode == "certificate" {
			req.Spec.Authentication.Certificate.CertData = []byte(req.CertDataStr)
			req.Spec.Authentication.Certificate.KeyData = []byte(req.KeyDataStr)
		}
		if req.Spec.Connect.Forward.ApiServer != "" {
			if !strings.HasPrefix(req.Spec.Connect.Forward.ApiServer, "https://") && !strings.HasPrefix(req.Spec.Connect.Forward.ApiServer, "http://") {
				req.Spec.Connect.Forward.ApiServer = fmt.Sprintf("%s%s", "https://", req.Spec.Connect.Forward.ApiServer)
			}
		}

		// 生成一个rsa格式的私钥
		privateKey, err := certificate.GeneratePrivateKey()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return

		}
		req.PrivateKey = privateKey
		client := kubernetes.NewKubernetes(&req.Cluster)
		if err := client.Ping(); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		v, _ := client.Version()
		req.Status.Version = v.GitVersion
		if req.Spec.Authentication.Mode == "configFile" {
			kubeCfg, err := client.Config()
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			req.Spec.Connect.Forward.ApiServer = kubeCfg.Host
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.CreatedBy = profile.Name

		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}
		req.Cluster.Status.Phase = clusterStatusSaved
		if err := h.clusterService.Create(&req.Cluster, txOptions); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		requiredPermissions := map[string][]string{
			"namespaces":       {"get", "post", "delete"},
			"clusterroles":     {"get", "post", "delete"},
			"clusterrolebings": {"get", "post", "delete"},
			"roles":            {"get", "post", "delete"},
			"rolebindings":     {"get", "post", "delete"},
		}
		notAllowed, err := checkRequiredPermissions(client, requiredPermissions)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if notAllowed != "" {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", []string{"permission %s required", notAllowed})
			return
		}
		_ = tx.Commit()
		ctx.Values().Set("data", &req)
		go func() {
			req.Status.Phase = clusterStatusInitializing
			if e := h.clusterService.Update(req.Name, &req.Cluster, common.DBOptions{}); e != nil {
				server.Logger().Errorf("can not update cluster status %s", err)
				return
			}
			if err := client.CreateDefaultClusterRoles(); err != nil {
				req.Status.Phase = clusterStatusFailed
				req.Status.Message = err.Error()
				if e := h.clusterService.Update(req.Name, &req.Cluster, common.DBOptions{}); e != nil {
					server.Logger().Errorf("can not update cluster status %s", err)
					return
				}
				server.Logger().Errorf("can not init  built in clusterroles %s", err)
				return
			}
			if !profile.IsAdministrator {
				binding := v1Cluster.Binding{
					BaseModel: v1.BaseModel{
						Kind: "ClusterBinding",
					},
					Metadata: v1.Metadata{
						Name: fmt.Sprintf("%s-%s-cluster-binding", req.Name, profile.Name),
					},
					UserRef:    profile.Name,
					ClusterRef: req.Name,
				}
				if err := h.clusterBindingService.CreateClusterBinding(&binding, common.DBOptions{}); err != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
				if err := client.CreateOrUpdateClusterRoleBinding("cluster-owner", profile.Name, true); err != nil {
					_ = tx.Rollback()
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
				if err := h.updateUserCert(client, &binding); err != nil {
					req.Status.Phase = clusterStatusFailed
					req.Status.Message = err.Error()
					if e := h.clusterService.Update(req.Name, &req.Cluster, common.DBOptions{}); e != nil {
						server.Logger().Errorf("can not update cluster status %s", err)
						return
					}
					server.Logger().Errorf("can not create cluster user  %s", err)
					return
				}
				req.Status.Phase = clusterStatusCompleted
				if e := h.clusterService.Update(req.Name, &req.Cluster, common.DBOptions{}); e != nil {
					server.Logger().Errorf("can not update cluster status %s", err)
					return
				}
			} else {
				req.Status.Phase = clusterStatusCompleted
				if e := h.clusterService.Update(req.Name, &req.Cluster, common.DBOptions{}); e != nil {
					server.Logger().Errorf("can not update cluster status %s", err)
					return
				}
			}
			err = client.CreateAppMarketCRD()
			if err != nil {
				server.Logger().Errorf("create app-market crd failed %s", err)
			}
		}()
	}
}

func (h *Handler) updateUserCert(client kubernetes.Interface, binding *v1Cluster.Binding) error {
	csr, err := client.CreateCommonUser(binding.UserRef)
	if err != nil {
		return err
	}
	binding.Certificate = csr
	if err := h.clusterBindingService.UpdateClusterBinding(binding.Name, binding, common.DBOptions{}); err != nil {
		return err
	}
	return nil
}

func checkRequiredPermissions(client kubernetes.Interface, requiredPermissions map[string][]string) (string, error) {
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	resultCh := make(chan kubernetes.PermissionCheckResult)
	doneCh := make(chan struct{})
	for key := range requiredPermissions {
		for i := range requiredPermissions[key] {
			wg.Add(1)
			i := i
			go func(key string, index int) {
				rs, err := client.HasPermission(authV1.ResourceAttributes{
					Verb:     requiredPermissions[key][i],
					Resource: key,
				})
				if err != nil {
					errCh <- err
					wg.Done()
					return
				}
				resultCh <- rs
				wg.Done()
			}(key, i)
		}
	}
	go func() {
		wg.Wait()
		doneCh <- struct{}{}
	}()
	for {
		select {
		case <-doneCh:
			goto end
		case err := <-errCh:
			return "", err
		case b := <-resultCh:
			if !b.Allowed {
				return fmt.Sprintf("%s", b.Resource.Resource), nil
			}
		}
	}
end:
	return "", nil
}

func (h *Handler) SearchClusters() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		showExtra := ctx.URLParamExists("showExtra")
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		clusters, total, err := h.clusterService.Search(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		result := make([]Cluster, 0)
		for i := range clusters {

			c := Cluster{Cluster: clusters[i]}
			if profile.IsAdministrator {
				c.Accessable = true
			} else {
				bs, err := h.clusterBindingService.GetClusterBindingByClusterName(c.Name, common.DBOptions{})
				if err != nil && !errors.Is(err, storm.ErrNotFound) {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
				c.MemberCount = len(bs)
				for j := range bs {
					if bs[j].UserRef == profile.Name {
						c.Accessable = true
					}
				}
			}
			result = append(result, c)
		}
		if showExtra {
			wg := &sync.WaitGroup{}
			for i := range result {
				c := kubernetes.NewKubernetes(&result[i].Cluster)
				wg.Add(1)
				i := i
				goCtx, cancel := goContext.WithTimeout(goContext.Background(), 1*time.Second)
				go func() {
					info, _ := getExtraClusterInfo(goCtx, c)
					result[i].ExtraClusterInfo = info
					cancel()
					wg.Done()
				}()
			}
			wg.Wait()
		}
		ctx.Values().Set("data", pkgV1.Page{Items: result, Total: total})
	}
}
func getExtraClusterInfo(context goContext.Context, client kubernetes.Interface) (ExtraClusterInfo, error) {
	err := client.Ping()
	if err != nil {
		return ExtraClusterInfo{Health: false}, nil
	}
	c, err := client.Client()
	if err != nil {
		return ExtraClusterInfo{}, err
	}
	nodesList, err := c.CoreV1().Nodes().List(context, metav1.ListOptions{})
	if err != nil {
		return ExtraClusterInfo{}, err
	}
	nodes := nodesList.Items

	totalCpu := float64(0)
	totalMemory := float64(0)
	usedCpu := float64(0)
	usedMemory := float64(0)
	readyNodes := 0
	for i := range nodes {
		conditions := nodes[i].Status.Conditions
		for i := range conditions {
			if conditions[i].Type == "Ready" {
				if conditions[i].Status == "True" {
					readyNodes += 1
				}
			}
		}
		cpu := nodes[i].Status.Allocatable.Cpu().AsApproximateFloat64()
		totalCpu += cpu
		memory := nodes[i].Status.Allocatable.Memory().AsApproximateFloat64()
		totalMemory += memory
	}
	podsList, err := c.CoreV1().Pods("").List(goContext.TODO(), metav1.ListOptions{})
	if err != nil {
		return ExtraClusterInfo{}, err
	}
	pods := podsList.Items
	for i := range pods {
		for j := range pods[i].Spec.Containers {
			cpu := pods[i].Spec.Containers[j].Resources.Requests.Cpu().AsApproximateFloat64()
			usedCpu += cpu
			memory := pods[i].Spec.Containers[j].Resources.Requests.Memory().AsApproximateFloat64()
			usedMemory += memory

		}
	}
	result := ExtraClusterInfo{
		Health:            true,
		TotalNodeNum:      len(nodes),
		ReadyNodeNum:      readyNodes,
		CPUAllocatable:    totalCpu,
		CPURequested:      usedCpu,
		MemoryAllocatable: totalMemory,
		MemoryRequested:   usedMemory,
	}
	return result, nil

}

func (h *Handler) GetCluster() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("get clusters failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", c)
	}
}

func (h *Handler) UpdateCluster() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req UpdateCluster
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		c.Labels = req.Labels
		if req.Labels == nil {
			req.Labels = []string{}
		}
		err = h.clusterService.Update(name, c, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) ListClusters() iris.Handler {
	return func(ctx *context.Context) {
		var clusters []v1Cluster.Cluster
		clusters, err := h.clusterService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		var resultClusters []Cluster
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		for i := range clusters {
			mbs, err := h.clusterBindingService.GetClusterBindingByClusterName(clusters[i].Name, common.DBOptions{})
			if err != nil && !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			rc := Cluster{
				Cluster: clusters[i],
			}
			for j := range mbs {
				if mbs[j].UserRef == profile.Name {
					rc.Accessable = true
				}
			}
			resultClusters = append(resultClusters, rc)
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", resultClusters)
	}
}

func (h *Handler) DeleteCluster() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		txOptions := common.DBOptions{DB: tx}

		if err := h.clusterService.Delete(name, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		clusterBindings, err := h.clusterBindingService.GetClusterBindingByClusterName(name, txOptions)
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		for i := range clusterBindings {
			if err := h.clusterBindingService.Delete(clusterBindings[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
				return
			}
		}
		k := kubernetes.NewKubernetes(c)
		_ = k.CleanAllRBACResource()
		_ = tx.Commit()
		ctx.StatusCode(iris.StatusOK)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/clusters")
	sp.Post("", handler.CreateCluster())
	sp.Get("", handler.ListClusters())
	sp.Get("/:name", handler.GetCluster())
	sp.Put("/:name", handler.UpdateCluster())
	sp.Delete("/:name", handler.DeleteCluster())
	sp.Post("/search", handler.SearchClusters())
	sp.Get("/:name/members", handler.ListClusterMembers())
	sp.Post("/:name/members", handler.CreateClusterMember())
	sp.Delete("/:name/members/:member", handler.DeleteClusterMember())
	sp.Put("/:name/members/:member", handler.UpdateClusterMember())
	sp.Get("/:name/members/:member", handler.GetClusterMember())
	sp.Get("/:name/clusterroles", handler.ListClusterRoles())
	sp.Post("/:name/clusterroles", handler.CreateClusterRole())
	sp.Put("/:name/clusterroles/:clusterrole", handler.UpdateClusterRole())
	sp.Delete("/:name/clusterroles/:clusterrole", handler.DeleteClusterRole())
	sp.Get("/:name/apigroups", handler.ListApiGroups())
	sp.Get("/:name/apigroups/{group:path}", handler.ListApiGroupResources())
	sp.Get("/:name/namespaces", handler.ListNamespace())
	sp.Get("/:name/terminal/session", handler.TerminalSessionHandler())
	sp.Get("/:name/logging/session", handler.LoggingHandler())
	sp.Get("/:name/repos", handler.ListClusterRepos())
	sp.Post("/:name/repos", handler.AddCLusterRepo())
	sp.Delete("/:name/repos/:repo", handler.DeleteClusterRepo())
}
