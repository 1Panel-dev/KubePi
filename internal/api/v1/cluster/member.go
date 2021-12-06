package cluster

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Cluster "github.com/KubeOperator/kubepi/internal/model/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/collectons"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

func (h *Handler) UpdateClusterMember() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req Member
		err := ctx.ReadJSON(&req)
		if err != nil {
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
		if c.CreatedBy == req.Name {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("can not delete or update cluster importer %s", req.Name))
			return
		}
		k := kubernetes.NewKubernetes(c)
		if err := k.CleanManagedClusterRoleBinding(req.Name); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		if err := k.CleanManagedRoleBinding(req.Name); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		// 删除重建
		for i := range req.NamespaceRoles {
			for j := range req.NamespaceRoles[i].Roles {
				if err := k.CreateOrUpdateRolebinding(req.NamespaceRoles[i].Namespace, req.NamespaceRoles[i].Roles[j], req.Name, false); err != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err)
					return
				}
			}
		}
		for i := range req.ClusterRoles {
			if err := k.CreateOrUpdateClusterRoleBinding(req.ClusterRoles[i], req.Name, false); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err)
				return
			}
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) GetClusterMember() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		memberName := ctx.Params().Get("member")

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
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}

		binding, err := h.clusterBindingService.GetBindingByClusterNameAndUserName(name, memberName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster binding failed: %s", err.Error()))
			return
		}
		labels := []string{
			fmt.Sprintf("%s=%s", kubernetes.LabelManageKey, "kubepi"),
			fmt.Sprintf("%s=%s", kubernetes.LabelClusterId, c.UUID),
			fmt.Sprintf("%s=%s", kubernetes.LabelUsername, binding.UserRef),
		}
		clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		rolebindings, err := client.RbacV1().RoleBindings("").List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: strings.Join(labels, ","),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}

		var member Member
		member.ClusterRoles = make([]string, 0)
		member.NamespaceRoles = make([]NamespaceRoles, 0)
		member.Name = binding.UserRef
		set := collectons.NewStringSet()
		for i := range clusterRoleBindings.Items {
			set.Add(clusterRoleBindings.Items[i].RoleRef.Name)
		}
		member.ClusterRoles = set.ToSlice()

		roleMap := map[string][]string{}

		for i := range rolebindings.Items {
			if roleMap[rolebindings.Items[i].Namespace] == nil {
				roleMap[rolebindings.Items[i].Namespace] = []string{rolebindings.Items[i].RoleRef.Name}
			} else {
				roleMap[rolebindings.Items[i].Namespace] = append(roleMap[rolebindings.Items[i].Namespace], rolebindings.Items[i].RoleRef.Name)
			}
		}
		for k := range roleMap {
			member.NamespaceRoles = append(member.NamespaceRoles, NamespaceRoles{
				Namespace: k,
				Roles:     roleMap[k],
			})
		}
		ctx.Values().Set("data", &member)
	}

}

func (h *Handler) ListClusterMembers() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		bindings, err := h.clusterBindingService.GetClusterBindingByClusterName(name, common.DBOptions{})
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		members := make([]Member, 0)
		for i := range bindings {
			members = append(members, Member{
				Name:        bindings[i].UserRef,
				BindingName: bindings[i].Name,
				CreateAt:    bindings[i].CreateAt,
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
			ctx.Values().Set("message", err.Error())
			return
		}
		if req.Name == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("username can not be none"))
			return
		}
		if len(req.ClusterRoles) == 0 && len(req.NamespaceRoles) == 0 {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("must select one role"))
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
			UserRef:    req.Name,
			ClusterRef: name,
		}

		tx, err := server.DB().Begin(true)
		c, err := h.clusterService.Get(name, common.DBOptions{DB: tx})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}

		k := kubernetes.NewKubernetes(c)
		cert, err := k.CreateCommonUser(req.Name)
		if err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("create common user failed: %s", err.Error()))
			return
		}
		binding.Certificate = cert
		if err := h.clusterBindingService.CreateClusterBinding(&binding, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "unable to complete authorization")
			return
		}
		// 创建clusterrolebinding
		for i := range req.ClusterRoles {
			if err := k.CreateOrUpdateClusterRoleBinding(req.ClusterRoles[i], req.Name, false); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", "unable to complete authorization")
				return
			}
		}
		// 创建Rolebinding
		for i := range req.NamespaceRoles {
			for j := range req.NamespaceRoles[i].Roles {
				if err := k.CreateOrUpdateRolebinding(req.NamespaceRoles[i].Namespace, req.NamespaceRoles[i].Roles[j], req.Name, false); err != nil {
					_ = tx.Rollback()
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err.Error())
					return
				}
			}
		}
		_ = tx.Commit()
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteClusterMember() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		memberName := ctx.Params().GetString("member")
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		if c.CreatedBy == memberName {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("can not delete or update cluster importer %s", profile.Name))
			return
		}

		binding, err := h.clusterBindingService.GetBindingByClusterNameAndUserName(c.Name, memberName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		if err := h.clusterBindingService.Delete(binding.Name, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster binding failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
		if err := k.CleanManagedClusterRoleBinding(memberName); err != nil {
			server.Logger().Errorf("can not delete cluster member %s : %s", memberName, err)
		}
		if err := k.CleanManagedRoleBinding(memberName); err != nil {
			server.Logger().Errorf("can not delete cluster member %s : %s", memberName, err)
		}
		_ = tx.Commit()
	}
}
