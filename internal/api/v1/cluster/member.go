package cluster

import (
	goContext "context"
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/collectons"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
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
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}
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
				b := rbacV1.RoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: req.NamespaceRoles[i].Namespace,
						Name:      fmt.Sprintf("%s-%s-%s", name, req.Name, req.NamespaceRoles[i].Roles[j]),
						Labels: map[string]string{
							kubernetes.LabelManageKey: "ekko",
							kubernetes.LabelClusterId: c.UUID,
							kubernetes.LabelUsername:              req.Name,
						},
						Annotations: map[string]string{
							"builtin":    "false",
							"created-at": time.Now().Format("2006-01-02 15:04:05"),
						},
					},
					Subjects: []rbacV1.Subject{
						{
							Kind: "User",
							Name: req.Name,
						},
					},
					RoleRef: rbacV1.RoleRef{
						Name: req.NamespaceRoles[i].Roles[j],
						Kind: "ClusterRole",
					},
				}
				_, err := client.RbacV1().RoleBindings(req.NamespaceRoles[i].Namespace).Create(goContext.TODO(), &b, metav1.CreateOptions{})
				if err != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", err)
					return
				}
			}
		}
		for i := range req.ClusterRoles {
			b := rbacV1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s-%s", name, req.Name, req.ClusterRoles[i]),
					Labels: map[string]string{
						kubernetes.LabelManageKey: "ekko",
						kubernetes.LabelClusterId: c.UUID,
						kubernetes.LabelUsername:               req.Name,
					},
					Annotations: map[string]string{
						"builtin":    "false",
						"created-at": time.Now().Format("2006-01-02 15:04:05"),
					},
				},
				Subjects: []rbacV1.Subject{
					{
						Kind: "User",
						Name: req.Name,
					},
				},
				RoleRef: rbacV1.RoleRef{
					Name: req.ClusterRoles[i],
					Kind: "ClusterRole",
				},
			}
			_, err := client.RbacV1().ClusterRoleBindings().Create(goContext.TODO(), &b, metav1.CreateOptions{})
			if err != nil {
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
			fmt.Sprintf("%s=%s", kubernetes.LabelManageKey, "ekko"),
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
		rolebindings, err := client.RbacV1().RoleBindings("").List(goContext.TODO(), metav1.ListOptions{LabelSelector: fmt.Sprintf("user-name=%s", binding.UserRef)})
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
		var members []Member
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
			ctx.Values().Set("message", fmt.Sprintf("create cluster member failed: %s", err.Error()))
			return
		}
		if req.Name == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Errorf("invalid username %s", req.Name))
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
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(c)
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
		// 创建clusterrolebinding
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}
		for i := range req.ClusterRoles {
			binding := rbacV1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s", req.Name, req.ClusterRoles[i]),
					Labels: map[string]string{
						kubernetes.LabelManageKey: "ekko",
						kubernetes.LabelClusterId: c.UUID,
						kubernetes.LabelUsername:               req.Name,
					},
					Annotations: map[string]string{
						"builtin":    "false",
						"created-at": time.Now().Format("2006-01-02 15:04:05"),
					},
				},
				Subjects: []rbacV1.Subject{
					{
						Kind: "User",
						Name: req.Name,
					},
				},
				RoleRef: rbacV1.RoleRef{
					Name: req.ClusterRoles[i],
					Kind: "ClusterRole",
				},
			}
			_, err := client.RbacV1().ClusterRoleBindings().Create(goContext.TODO(), &binding, metav1.CreateOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("create cluster role binding failed: %s", err.Error()))
				return
			}
		}
		// 创建Rolebinding
		for i := range req.NamespaceRoles {
			for j := range req.NamespaceRoles[i].Roles {
				b := rbacV1.RoleBinding{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: req.NamespaceRoles[i].Namespace,
						Name:      fmt.Sprintf("%s-%s", req.Name, req.NamespaceRoles[i].Roles[j]),
						Labels: map[string]string{
							"kubeoperator.io/manage": "ekko",
							kubernetes.LabelUsername:              req.Name,
						},
						Annotations: map[string]string{
							"builtin":    "false",
							"created-at": time.Now().Format("2006-01-02 15:04:05"),
						},
					},
					Subjects: []rbacV1.Subject{
						{
							Kind: "User",
							Name: req.Name,
						},
					},
					RoleRef: rbacV1.RoleRef{
						Kind: "ClusterRole",
						Name: req.NamespaceRoles[i].Roles[j],
					},
				}
				_, err := client.RbacV1().RoleBindings(req.NamespaceRoles[i].Namespace).Create(goContext.TODO(), &b, metav1.CreateOptions{})
				if err != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", fmt.Sprintf("create role binding failed: %s", err.Error()))
					return
				}
			}
		}

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
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		if err := k.CleanManagedRoleBinding(memberName); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		_ = tx.Commit()
	}
}
