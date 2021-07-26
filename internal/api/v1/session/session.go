package session

import (
	goContext "context"
	"errors"
	"fmt"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/internal/service/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/rolebinding"
	"github.com/KubeOperator/ekko/internal/service/v1/user"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/KubeOperator/ekko/pkg/collectons"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/crypto/bcrypt"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handler struct {
	userService        user.Service
	roleService        role.Service
	clusterService     cluster.Service
	rolebindingService rolebinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		clusterService:     cluster.NewService(),
		userService:        user.NewService(),
		roleService:        role.NewService(),
		rolebindingService: rolebinding.NewService(),
	}
}

func (h *Handler) IsLogin() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", loginUser != nil)
	}
}

func (h *Handler) Login() iris.Handler {
	return func(ctx *context.Context) {
		var loginCredential LoginCredential
		if err := ctx.ReadJSON(&loginCredential); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u, err := h.userService.GetByNameOrEmail(loginCredential.Username, common.DBOptions{})
		if err != nil {
			if errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", fmt.Sprintf("user %s: not Found", loginCredential.Username))
				return
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("query user %s failed ,: %s", loginCredential.Username, err.Error()))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Authenticate.Password), []byte(loginCredential.Password)); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "username or password error")
			return
		}

		permissions, err := h.aggregateResourcePermissions(loginCredential.Username)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		session := sessions.Get(ctx)
		profile := UserProfile{
			Name:                u.Name,
			NickName:            u.NickName,
			Email:               u.Email,
			ResourcePermissions: permissions,
		}
		session.Set("profile", profile)
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", profile)
	}
}

func (h *Handler) aggregateResourcePermissions(name string) (map[string][]string, error) {
	userRoleBindings, err := h.rolebindingService.GetRoleBindingBySubject(v1Role.Subject{
		Kind: "User",
		Name: name,
	}, common.DBOptions{})
	if err != nil && !errors.As(err, &storm.ErrNotFound) {
		return nil, err
	}

	allRoleBindings := append(userRoleBindings)
	var roleNames []string
	for i := range allRoleBindings {
		roleNames = append(roleNames, allRoleBindings[i].RoleRef)
	}
	rs, _, err := h.roleService.Search(0, 0, pkgV1.Conditions{
		"Name": pkgV1.Condition{
			Field:    "Name",
			Operator: "in",
			Value:    roleNames,
		},
	}, common.DBOptions{})
	if err != nil && !errors.As(err, &storm.ErrNotFound) {
		return nil, err
	}
	mapping := map[string]*collectons.StringSet{}
	var policyRoles []v1Role.PolicyRule
	//merge permissions
	for i := range rs {
		for j := range rs[i].Rules {
			policyRoles = append(policyRoles, rs[i].Rules[j])
		}
	}
	for i := range policyRoles {
		for j := range policyRoles[i].Resource {
			_, ok := mapping[policyRoles[i].Resource[j]]
			if !ok {
				mapping[policyRoles[i].Resource[j]] = collectons.NewStringSet()
			}
			for k := range policyRoles[i].Verbs {
				mapping[policyRoles[i].Resource[j]].Add(policyRoles[i].Verbs[k])
				if policyRoles[i].Verbs[k] == "*" {
					break
				}
			}
			if policyRoles[i].Resource[j] == "*" {
				break
			}
		}
	}
	resourceMapping := map[string][]string{}
	for key := range mapping {
		resourceMapping[key] = mapping[key].ToSlice()
	}
	return resourceMapping, nil
}

func (h *Handler) Logout() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		session.Delete("profile")
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", "logout success")
	}
}

func (h *Handler) GetProfile() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		permissions, err := h.aggregateResourcePermissions(p.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		p.ResourcePermissions = permissions
		session.Set("profile", p)
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", p)
	}
}

func (h *Handler) ListUserNamespace() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("cluster_name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		session := sessions.Get(ctx)
		u := session.Get("profile")
		profile := u.(UserProfile)

		k := kubernetes.NewKubernetes(*c)
		ns, err := k.GetUserNamespaceNames(profile.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		ctx.Values().Set("data", ns)
	}
}

func (h *Handler) GetClusterProfile() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		clusterName := ctx.Params().GetString("cluster_name")
		c, err := h.clusterService.Get(clusterName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		k := kubernetes.NewKubernetes(*c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get k8s client failed: %s", err.Error()))
			return
		}
		u := session.Get("profile")
		profile := u.(UserProfile)
		clusterRoleBindings, err := client.RbacV1().ClusterRoleBindings().List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("user-name=%s", profile.Name),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster-role-binding failed: %s", err.Error()))
			return
		}
		rolebindings, err := client.RbacV1().RoleBindings("").List(goContext.TODO(), metav1.ListOptions{
			LabelSelector: fmt.Sprintf("user-name=%s", profile.Name),
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get role-binding failed: %s", err.Error()))
			return
		}
		roleSet := map[string]struct{}{}
		for i := range clusterRoleBindings.Items {
			for j := range clusterRoleBindings.Items[i].Subjects {
				if clusterRoleBindings.Items[i].Subjects[j].Kind == "User" {
					roleSet[clusterRoleBindings.Items[i].RoleRef.Name] = struct{}{}
				}
			}
		}
		for i := range rolebindings.Items {
			for j := range rolebindings.Items[i].Subjects {
				if rolebindings.Items[i].Subjects[j].Kind == "User" {
					roleSet[rolebindings.Items[i].RoleRef.Name] = struct{}{}
				}
			}
		}
		var roles []v1.ClusterRole
		for key := range roleSet {
			r, err := client.RbacV1().ClusterRoles().Get(goContext.TODO(), key, metav1.GetOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("get cluster-role failed: %s", err.Error()))
				return
			}
			roles = append(roles, *r)
		}
		crp := ClusterUserProfile{
			UserProfile:  profile,
			ClusterRoles: roles,
		}
		ctx.Values().Set("data", &crp)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/sessions")
	sp.Post("", handler.Login())
	sp.Delete("", handler.Logout())
	sp.Get("", handler.GetProfile())
	sp.Get("/:cluster_name", handler.GetClusterProfile())
	sp.Get("/status", handler.IsLogin())
	sp.Get("/:cluster_name/namespaces", handler.ListUserNamespace())
}
