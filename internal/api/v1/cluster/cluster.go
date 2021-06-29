package cluster

import (
	goContext "context"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/KubeOperator/ekko/pkg/certificate"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	authV1 "k8s.io/api/authorization/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sync"
)

type Handler struct {
	clusterService        cluster.Service
	clusterBindingService clusterbinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		clusterService:        cluster.NewService(),
		clusterBindingService: clusterbinding.NewService(),
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
		if req.CaDataStr != "" {
			req.CaCertificate.CertData = []byte(req.CaDataStr)
		}
		if req.Spec.Authentication.Mode == "certificate" {
			req.Spec.Authentication.Certificate.CertData = []byte(req.CertDataStr)
			req.Spec.Authentication.Certificate.KeyData = []byte(req.KeyDataStr)
		}
		// 生成一个rsa格式的私钥
		privateKey, err := certificate.GeneratePrivateKey()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return

		}
		req.PrivateKey = privateKey

		client := kubernetes.NewKubernetes(req.Cluster)
		if err := client.Ping(); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		v, _ := client.Version()
		req.Status.Version = v.GitVersion
		requiredPermissions := map[string][]string{
			"namespaces":       {"get", "post", "delete"},
			"clusterroles":     {"get", "post", "delete"},
			"clusterrolebings": {"get", "post", "delete"},
			"roles":            {"get", "post", "delete"},
			"rolebindings":     {"get", "post", "delete"},
		}
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
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			case b := <-resultCh:
				if !b.Allowed {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", fmt.Errorf("permission %s-%s  not allowed", b.Resource.Resource, b.Resource.Verb))
					return
				}
			}
		}
	end:
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)

		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}

		req.CreatedBy = profile.Name
		if err := client.CreateDefaultClusterRoles(); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.clusterService.Create(&req.Cluster, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		// 把创建者以管理员身份加入集群中

		binding := v1Cluster.Binding{
			BaseModel: v1.BaseModel{
				Kind:      "ClusterBinding",
				CreatedBy: profile.Name,
			},
			Metadata: v1.Metadata{
				Name: fmt.Sprintf("%s-%s-cluster-binding-", req.Name, profile.Name),
			},
			Subject: v1Cluster.Subject{
				Name: profile.Name,
				Kind: "User",
			},
			ClusterRef: req.Name,
		}

		csr, err := client.CreateCommonUser(profile.Name)
		if err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		binding.Certificate = csr

		if err := h.clusterBindingService.CreateClusterBinding(&binding, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		// 管理员权限

		kc, err := client.Client()
		if err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		if _, err := kc.RbacV1().ClusterRoleBindings().Create(goContext.TODO(), &rbacV1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("ekko:cluster-admin-%s", profile.Name),
			},
			Subjects: []rbacV1.Subject{
				{
					Kind: "User",
					Name: profile.Name,
				},
			},
			RoleRef: rbacV1.RoleRef{
				Name: "cluster-admin",
				Kind: "ClusterRole",
			},
		}, metav1.CreateOptions{}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()

		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) SearchClusters() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions pkgV1.Conditions
		if ctx.GetContentLength() > 0 {
			if err := ctx.ReadJSON(&conditions); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		clusters, total, err := h.clusterService.Search(pageNum, pageSize, conditions, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", pkgV1.Page{Items: clusters, Total: total})
	}
}

func (h *Handler) ListClusters() iris.Handler {
	return func(ctx *context.Context) {
		var clusters []v1Cluster.Cluster
		clusters, err := h.clusterService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("get clusters failed: %s", err.Error()))
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", clusters)
	}
}

func (h *Handler) DeleteCluster() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")

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
		if err != nil {
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
		_ = tx.Commit()
		ctx.StatusCode(iris.StatusOK)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/clusters")
	sp.Post("", handler.CreateCluster())
	sp.Get("", handler.ListClusters())
	sp.Delete("/:name", handler.DeleteCluster())
	sp.Post("/search", handler.SearchClusters())
	sp.Get("/:name/members", handler.GetClusterMembers())
	sp.Post("/:name/members", handler.CreateClusterMember())
	sp.Delete("/:name/members/:member", handler.DeleteClusterMember())
	sp.Put("/:name/members/:member", handler.UpdateClusterMember())
	sp.Get("/:name/members/:member", handler.GetClusterMember())
	sp.Get("/:name/clusterroles", handler.GetClusterRoles())
	sp.Post("/:name/clusterroles", handler.CreateClusterRole())
	sp.Put("/:name/clusterroles/:clusterrole", handler.UpdateClusterRole())
	sp.Delete("/:name/clusterroles/:clusterrole", handler.DeleteClusterRole())
	sp.Get("/:name/apigroups", handler.ListApiGroups())
	sp.Get("/:name/apigroups/{group:path}", handler.ListApiGroupResources())

}
