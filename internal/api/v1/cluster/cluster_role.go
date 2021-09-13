package cluster

import (
	goContext "context"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
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
		for i := range req.Rules {
			for j := range req.Rules[i].APIGroups {
				if req.Rules[i].APIGroups[j] == "core" {
					req.Rules[i].APIGroups[j] = ""
				}
			}
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
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
		instance.Annotations["description"] = req.Annotations["description"]
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

		if req.Name == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("username can not be none"))
			return
		}
		if len(req.Rules) == 0 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("must select one rule"))
			return
		}

		for i := range req.Rules {
			for j := range req.Rules[i].APIGroups {
				if req.Rules[i].APIGroups[j] == "core" {
					req.Rules[i].APIGroups[j] = ""
				}
			}
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		req.Annotations["builtin"] = "false"
		req.Annotations["created-at"] = time.Now().Format("2006-01-02 15:04:05")
		req.Labels[kubernetes.LabelManageKey] = "kubepi"
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
		k := kubernetes.NewKubernetes(c)
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

func (h *Handler) ListClusterRoles() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		scope := ""
		if ctx.URLParamExists("scope") {
			scope = ctx.URLParam("scope")
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get kubernetes client failed: %s", err.Error()))
			return
		}
		labels := []string{
			fmt.Sprintf("%s=%s", kubernetes.LabelManageKey, "kubepi"),
		}
		if scope != "" {
			labels = append(labels, fmt.Sprintf("%s=%s", "kubepi.org/role-type", scope))
		}

		items, err := client.RbacV1().ClusterRoles().List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster roles failed: %s", err.Error()))
			return
		}
		ctx.Values().Set("data", items.Items)
	}
}
