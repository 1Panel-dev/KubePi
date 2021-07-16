package session

import (
	"errors"
	"fmt"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/internal/service/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/rolebinding"
	"github.com/KubeOperator/ekko/internal/service/v1/user"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/KubeOperator/ekko/pkg/collectons"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userService         user.Service
	roleService         role.Service
	rolebindingService  rolebinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService:         user.NewService(),
		roleService:         role.NewService(),
		rolebindingService:  rolebinding.NewService(),
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

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/sessions")
	sp.Post("", handler.Login())
	sp.Delete("", handler.Logout())
	sp.Get("", handler.GetProfile())
	sp.Get("/status", handler.IsLogin())
}
