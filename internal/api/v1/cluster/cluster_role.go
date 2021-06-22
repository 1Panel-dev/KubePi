package cluster

import (
	goContext "context"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func (h *Handler) UpdateClusterRole() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		clusterRoleName := ctx.Params().GetString("clusterrole")

		var req rbacV1.ClusterRole
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
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
		instance, err := client.RbacV1().ClusterRoles().Get(goContext.TODO(), clusterRoleName, metav1.GetOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		instance.Rules = req.Rules
		resp, err := client.RbacV1().ClusterRoles().Update(goContext.TODO(), instance, metav1.UpdateOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("create cluster role failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", resp)
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
		c, err := h.clusterService.Get(name, common.DBOptions{})
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
func (h *Handler) DeleteClusterRole() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		clusterRole := ctx.Params().GetString("clusterrole")
		c, err := h.clusterService.Get(name, common.DBOptions{})
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
		r, err := client.RbacV1().ClusterRoles().Get(goContext.TODO(), clusterRole, metav1.GetOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster role   failed: %s", err.Error()))
			return
		}
		createBy, ok := r.Annotations["created-by"]
		if ok {
			if createBy == "system" {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", fmt.Sprintf("can not delete it ,beacuse it created by system"))
				return
			}
		}
		if err := client.RbacV1().ClusterRoles().Delete(goContext.TODO(), r.Name, metav1.DeleteOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("delete  cluster role failed: %s", err.Error()))
			return
		}
	}
}

func (h *Handler) GetClusterRoles() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
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
