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

		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}

		binding, err := h.clusterBindingService.GetBindingByClusterNameAndSubject(name, v1Cluster.Subject{
			Kind: req.Kind,
			Name: req.Name,
		}, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster binding failed: %s", err.Error()))
			return
		}

		clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("subject-name=%s,subject-kind=%s", binding.Subject.Name, binding.Subject.Kind),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster role binding failed: %s", err.Error()))
			return
		}
		var forCreate []string
		var forDelete []string

		for i := range req.ClusterRoles {
			exists := false
			for j := range clusterRoleBindings.Items {
				if req.ClusterRoles[i] == clusterRoleBindings.Items[j].RoleRef.Name {
					exists = true
				}
			}
			if !exists {
				forCreate = append(forCreate, req.ClusterRoles[i])
			}
		}
		for i := range clusterRoleBindings.Items {
			exists := false
			for j := range req.ClusterRoles {
				if req.ClusterRoles[j] == clusterRoleBindings.Items[i].RoleRef.Name {
					exists = true
				}
			}
			if !exists {
				forDelete = append(forDelete, clusterRoleBindings.Items[i].Name)
			}
		}
		for i := range forDelete {
			if err := client.RbacV1().ClusterRoleBindings().Delete(goContext.TODO(), forDelete[i], metav1.DeleteOptions{}); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("update cluster role binding failed: %s", err.Error()))
				return
			}
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		for i := range forCreate {
			b := rbacV1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: fmt.Sprintf("%s-%s-%s", name, req.Name, forCreate[i]),
					Labels: map[string]string{
						"manage":       "ekko",
						"subject-name": req.Name,
						"subject-kind": req.Kind,
					},
					Annotations: map[string]string{
						"created-by": profile.Name,
						"created-at": time.Now().Format("2006-01-02 15:04:05"),
					},
				},
				Subjects: []rbacV1.Subject{
					{
						Kind: req.Kind,
						Name: req.Name,
					},
				},
				RoleRef: rbacV1.RoleRef{
					Name: forCreate[i],
					Kind: "ClusterRole",
				},
			}
			_, err := client.RbacV1().ClusterRoleBindings().Create(goContext.TODO(), &b, metav1.CreateOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("create cluster role binding failed: %s", err.Error()))
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

		var kind string
		if ctx.URLParamExists("kind") {
			kind = ctx.URLParam("kind")
		}
		subject := v1Cluster.Subject{
			Name: memberName,
		}
		if kind != "" {
			switch strings.ToLower(kind) {
			case "user":
				subject.Kind = "User"
			case "group":
				subject.Kind = "Group"
			default:
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", fmt.Errorf("invalid kind %s", kind)))
				return
			}
			subject.Kind = kind
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
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}

		binding, err := h.clusterBindingService.GetBindingByClusterNameAndSubject(name, subject, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster binding failed: %s", err.Error()))
			return
		}

		clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("subject-name=%s,subject-kind=%s", binding.Subject.Name, binding.Subject.Kind),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster role binding failed: %s", err.Error()))
			return
		}

		var member Member
		member.ClusterRoles = make([]string, 0)
		member.Kind = binding.Subject.Kind
		member.Name = binding.Subject.Name
		set := collectons.NewStringSet()
		for i := range clusterRoleBindings.Items {
			set.Add(clusterRoleBindings.Items[i].RoleRef.Name)
		}
		member.ClusterRoles = set.ToSlice()
		ctx.Values().Set("data", &member)
	}

}

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
				Name: fmt.Sprintf("%s-%s-cluster-binding-", name, req.Name),
			},
			Subject: v1Cluster.Subject{
				Name: req.Name,
				Kind: req.Kind,
			},
			ClusterRef: name,
		}
		// 生成用户证书
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		// 查询相关用户组

		gs, err := h.groupBindingService.ListByUserName(req.Name, common.DBOptions{})
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("list user failed: %s", err.Error()))
			return
		}
		groupNames := make([]string, 0)
		for i := range gs {
			groupNames = append(groupNames, gs[i].GroupRef)
		}
		cert, err := k.CreateCommonUser(req.Name, groupNames...)
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
					Name: fmt.Sprintf("%s-%s-%s", name, req.Name, req.ClusterRoles[i]),
					Labels: map[string]string{
						"manage":       "ekko",
						"subject-name": req.Name,
						"subject-kind": req.Kind,
					},
					Annotations: map[string]string{
						"created-by": profile.Name,
						"created-at": time.Now().Format("2006-01-02 15:04:05"),
					},
				},
				Subjects: []rbacV1.Subject{
					{
						Kind: req.Kind,
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
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteClusterMember() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		member := ctx.Params().GetString("member")
		var kind string
		if ctx.URLParamExists("kind") {
			kind = ctx.URLParam("kind")
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		subject := v1Cluster.Subject{
			Name: member,
		}
		if kind != "" {
			switch strings.ToLower(kind) {
			case "user":
				subject.Kind = "User"
			case "group":
				subject.Kind = "Group"
			default:
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", fmt.Errorf("invalid kind %s", kind)))
				return
			}
			subject.Kind = kind
		}

		binding, err := h.clusterBindingService.GetBindingByClusterNameAndSubject(c.Name, subject, common.DBOptions{})
		if err != nil {
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
		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}
		if err = client.RbacV1().ClusterRoleBindings().DeleteCollection(goContext.TODO(),
			metav1.DeleteOptions{},
			metav1.ListOptions{
				LabelSelector: fmt.Sprintf("subject-name=%s,subject-kind=%s", subject.Name, subject.Kind),
			}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster role failed: %s", err.Error()))
			return
		}
		_ = tx.Commit()
	}
}
