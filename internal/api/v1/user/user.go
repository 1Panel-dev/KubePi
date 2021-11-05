package user

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	v1User "github.com/KubeOperator/kubepi/internal/model/v1/user"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/rolebinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/user"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/KubeOperator/kubepi/pkg/collectons"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	userService           user.Service
	roleBindingService    rolebinding.Service
	clusterBindingService clusterbinding.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService:           user.NewService(),
		roleBindingService:    rolebinding.NewService(),
		clusterBindingService: clusterbinding.NewService(),
	}
}

func (h *Handler) SearchUsers() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)

		//pattern := ctx.URLParam("pattern")
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		users, total, err := h.userService.Search(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		us := make([]User, 0)
		for i := range users {
			bindings, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{Kind: "User", Name: users[i].Name}, common.DBOptions{})
			if err != nil && !errors.As(err, &storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			roles := collectons.NewStringSet()
			for i := range bindings {
				roles.Add(bindings[i].RoleRef)
			}
			us = append(us, User{
				User:  users[i],
				Roles: roles.ToSlice(),
			})
		}
		ctx.Values().Set("data", pkgV1.Page{Items: us, Total: total})
	}
}

func (h *Handler) CreateUser() iris.Handler {
	return func(ctx *context.Context) {
		var req User
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		//tx
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.CreatedBy = profile.Name
		if req.Language == "" {
			req.Language = profile.Language
		}
		req.Type = v1User.LOCAL
		if err := h.userService.Create(&req.User, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			if errors.Is(err, storm.ErrAlreadyExists) {
				u, _ := h.userService.GetByNameOrEmail(req.User.Name, common.DBOptions{})
				if u != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", "username already exists")
					return
				}
				u, _ = h.userService.GetByNameOrEmail(req.User.Email, common.DBOptions{})
				if u != nil {
					ctx.StatusCode(iris.StatusInternalServerError)
					ctx.Values().Set("message", "email already exists")
					return
				}
			}
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if len(req.Roles) > 0 {
			for i := range req.Roles {
				roleName := req.Roles[i]
				binding := v1Role.Binding{
					BaseModel: v1.BaseModel{
						Kind:       "RoleBind",
						ApiVersion: "v1",
						CreatedBy:  profile.Name,
					},
					Metadata: v1.Metadata{
						Name: fmt.Sprintf("role-binding-%s-%s", roleName, req.Name),
					},
					Subject: v1Role.Subject{
						Kind: "User",
						Name: req.Name,
					},
					RoleRef: roleName,
				}
				if err := h.roleBindingService.CreateRoleBinding(&binding, common.DBOptions{DB: tx}); err != nil {
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

func (h *Handler) DeleteUser() iris.Handler {
	return func(ctx *context.Context) {
		userName := ctx.Params().GetString("name")

		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)

		if userName == profile.Name {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Errorf("can not delete yourself"))
			return
		}
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}

		rbs, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{
			Kind: "User",
			Name: userName,
		}, txOptions)
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		for i := range rbs {
			if err := h.roleBindingService.Delete(rbs[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		cbs, err := h.clusterBindingService.GetBindingsByUserName(userName, txOptions)
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		for i := range cbs {
			if err := h.clusterBindingService.Delete(cbs[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		if err := h.userService.Delete(userName, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
	}
}

func (h *Handler) GetUser() iris.Handler {
	return func(ctx *context.Context) {
		userName := ctx.Params().GetString("name")
		u, err := h.userService.GetByNameOrEmail(userName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		bindings, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{Kind: "User", Name: u.Name}, common.DBOptions{})
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		roles := collectons.NewStringSet()
		for i := range bindings {
			roles.Add(bindings[i].RoleRef)
		}
		ctx.Values().Set("data", &User{User: *u, Roles: roles.ToSlice()})
	}
}

func (h *Handler) GetUsers() iris.Handler {
	return func(ctx *context.Context) {
		us, err := h.userService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", us)
	}
}

func (h *Handler) UpdateUser() iris.Handler {
	return func(ctx *context.Context) {
		userName := ctx.Params().GetString("name")
		var req User
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if req.Password != "" {
			if err := h.userService.ResetPassword(userName, req.Password, common.DBOptions{}); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			ctx.Values().Set("data", "ok")
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.userService.Update(userName, &req.User, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		bindings, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{Kind: "User", Name: userName}, common.DBOptions{})
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		currentRoles := collectons.NewStringSet()
		for i := range bindings {
			currentRoles.Add(bindings[i].RoleRef)
		}
		for i := range req.Roles {
			r := req.Roles[i]
			if currentRoles.Exists(r) {
				continue
			}
			binding := v1Role.Binding{
				BaseModel: v1.BaseModel{
					Kind:       "RoleBind",
					ApiVersion: "v1",
					CreatedBy:  profile.Name,
				},
				Metadata: v1.Metadata{
					Name: fmt.Sprintf("role-binding-%s-%s", r, req.Name),
				},
				Subject: v1Role.Subject{
					Kind: "User",
					Name: req.Name,
				},
				RoleRef: r,
			}
			if err := h.roleBindingService.CreateRoleBinding(&binding, common.DBOptions{DB: tx}); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			currentRoles.Add(binding.RoleRef)
		}
		diffs := currentRoles.Difference(req.Roles)

		for i := range bindings {
			for j := range diffs {
				if bindings[i].RoleRef == diffs[j] {
					if err := h.roleBindingService.Delete(bindings[i].Name, common.DBOptions{DB: tx}); err != nil {
						_ = tx.Rollback()
						ctx.StatusCode(iris.StatusInternalServerError)
						ctx.Values().Set("message", err.Error())
						return
					}
				}
			}
		}

		_ = tx.Commit()
		ctx.Values().Set("data", &req)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/users")
	sp.Post("/search", handler.SearchUsers())
	sp.Post("/", handler.CreateUser())
	sp.Delete("/:name", handler.DeleteUser())
	sp.Get("/:name", handler.GetUser())
	sp.Put("/:name", handler.UpdateUser())
	sp.Get("/", handler.GetUsers())
}
