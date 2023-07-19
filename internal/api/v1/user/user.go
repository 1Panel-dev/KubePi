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
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/rolebinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/user"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/KubeOperator/kubepi/pkg/collectons"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	userService           user.Service
	roleBindingService    rolebinding.Service
	clusterBindingService clusterbinding.Service
	clusterService        cluster.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService:           user.NewService(),
		roleBindingService:    rolebinding.NewService(),
		clusterBindingService: clusterbinding.NewService(),
		clusterService:        cluster.NewService(),
	}
}

// Search User
// @Tags users
// @Summary Search users
// @Description Search users by Condition
// @Accept  json
// @Produce  json
// @Success 200 {object} api.Page
// @Security ApiKeyAuth
// @Router /users/search [post]
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
			users[i].Authenticate = v1User.Authenticate{}
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

// Create User
// @Tags users
// @Summary Create user
// @Description Create user
// @Accept  json
// @Produce  json
// @Param request body docs.UserCreate true "request"
// @Success 200 {object} v1User.User
// @Security ApiKeyAuth
// @Router /users [post]
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
		req.Authenticate = v1User.Authenticate{}
		ctx.Values().Set("data", req)
	}
}

// Delete User
// @Tags users
// @Summary Delete user by name
// @Description Delete user by name
// @Accept  json
// @Produce  json
// @Param name path string true "用户名称"
// @Success 200 {object} v1User.User
// @Security ApiKeyAuth
// @Router /users/{name} [delete]
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
			c, err := h.clusterService.Get(cbs[i].ClusterRef, common.DBOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
				return
			}
			k := kubernetes.NewKubernetes(c)
			if err := k.CleanManagedClusterRoleBinding(cbs[i].UserRef); err != nil {
				server.Logger().Errorf("can not delete cluster member %s : %s", cbs[i].UserRef, err)
			}
			if err := k.CleanManagedRoleBinding(cbs[i].UserRef); err != nil {
				server.Logger().Errorf("can not delete cluster member %s : %s", cbs[i].UserRef, err)
			}

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

// Get User
// @Tags users
// @Summary Get user by name
// @Description Get user by name
// @Accept  json
// @Produce  json
// @Param name path string true "用户名称"
// @Success 200 {object} v1User.User
// @Security ApiKeyAuth
// @Router /users/{name} [get]
func (h *Handler) GetUser() iris.Handler {
	return func(ctx *context.Context) {
		userName := ctx.Params().GetString("name")
		u, err := h.userService.GetByNameOrEmail(userName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		u.Authenticate = v1User.Authenticate{}
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

// List User
// @Tags users
// @Summary List all users
// @Description List all users
// @Accept  json
// @Produce  json
// @Success 200 {object} []v1User.User
// @Security ApiKeyAuth
// @Router /users [get]
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

// Update User
// @Tags users
// @Summary Update user by name
// @Description Update user by name
// @Accept  json
// @Produce  json
// @Param name path string true "用户名称"
// @Success 200 {object} v1User.User
// @Security ApiKeyAuth
// @Router /users/{name} [put]
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
			if err := h.userService.UpdatePassword(userName, req.OldPassword, req.Password, common.DBOptions{}); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", "can not match original password")
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
