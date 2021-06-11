package role

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	"github.com/KubeOperator/ekko/internal/service/v1/role"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	roleService role.Service
}

func NewHandler() *Handler {
	return &Handler{
		roleService: role.NewService(),
	}
}

func (h *Handler) SearchRoles() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions pkgV1.Conditions
		if ctx.GetContentLength() > 0 {
			if err := ctx.ReadJSON(&conditions); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		groups, total, err := h.roleService.Search(pageNum, pageSize, conditions)
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: groups, Total: total})
	}
}

func (h *Handler) ListRoles() iris.Handler {
	return func(ctx *context.Context) {
		roles, err := h.roleService.List()
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", roles)
	}
}

func (h *Handler) CreateRole() iris.Handler {
	return func(ctx *context.Context) {
		var req Role
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.CreatedBy = profile.Name
		if req.TemplateRef != "" {
			if err := h.roleService.CreateWithTemplate(&req.Role, req.TemplateRef); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		} else {
			if err := h.roleService.Create(&req.Role); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", req.Role)
	}
}

func (h *Handler) DeleteRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		if roleName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", roleName))
			return
		}
		if err := h.roleService.Delete(roleName); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) UpdateRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		if roleName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", roleName))
			return
		}
		var req Role
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.roleService.Update(roleName, &req.Role); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) GetRole() iris.Handler {
	return func(ctx *context.Context) {
		roleName := ctx.Params().GetString("name")
		if roleName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", roleName))
			return
		}
		r, err := h.roleService.Get(roleName)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", r)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/roles")
	sp.Post("/search", handler.SearchRoles())
	sp.Get("/", handler.ListRoles())
	sp.Get("/:name", handler.GetRole())
	sp.Post("/", handler.CreateRole())
	sp.Delete("/:name", handler.DeleteRole())
	sp.Put("/:name", handler.UpdateRole())
}
