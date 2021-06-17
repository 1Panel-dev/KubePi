package cluster

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (h *Handler) GetClusterMembers() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		bindings, err := h.clusterBindingService.GetClusterBindingByClusterName(name, common.DBOptions{})
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		subjectMap := map[v1Cluster.Subject]v1Cluster.Binding{}
		for i := range bindings {
			subjectMap[bindings[i].Subject] = bindings[i]
		}
		members := make([]Member, 0)
		for key := range subjectMap {
			members = append(members, Member{
				Name:        key.Name,
				Kind:        key.Kind,
				BindingName: subjectMap[key].Name,
				CreateAt:    subjectMap[key].CreateAt,
			})
		}
		ctx.Values().Set("data", members)
	}
}

func (h *Handler) CreateClusterMember() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req Member
		err := ctx.ReadJSON(&req)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		binding := v1Cluster.Binding{
			BaseModel: v1.BaseModel{
				Kind:      "ClusterBinding",
				CreatedBy: profile.Name,
			},
			Metadata: v1.Metadata{
				Name: fmt.Sprintf("%s-%s-cluster-binding", name, req.Name),
			},
			Subject: v1Cluster.Subject{
				Name: req.Name,
				Kind: req.Kind,
			},
			ClusterRef: name,
		}
		// 生成用户证书
		c, err := h.clusterService.Get(name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		cert, err := k.CreateCommonUser(req.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("create common user failed: %s", err.Error()))
			return
		}
		binding.Certificate = cert
		if err := h.clusterBindingService.CreateClusterBinding(&binding, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) CreateClusterRole() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req rbacV1.ClusterRole
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		c, err := h.clusterService.Get(name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.Annotations = map[string]string{
			"ekko-i18n":  "none",
			"created-by": profile.Name,
			"created-at": time.Now().Format("2006-01-02 15:04:05"),
		}
		req.Labels = map[string]string{
			"manage": "ekko",
		}
		resp, err := client.RbacV1().ClusterRoles().Create(goContext.TODO(), &req, metav1.CreateOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("create cluster role failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", resp)
	}
}

func (h *Handler) GetClusterRoles() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get kubernetes client failed: %s", err.Error()))
			return
		}
		items, err := client.RbacV1().ClusterRoles().List(goContext.TODO(), metav1.ListOptions{LabelSelector: "manage=ekko"})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster roles failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", items.Items)
	}
}
