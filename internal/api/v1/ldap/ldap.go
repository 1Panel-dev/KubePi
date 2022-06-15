package ldap

import (
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/ldap"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	ldapService ldap.Service
}

func NewHandler() *Handler {
	return &Handler{
		ldapService: ldap.NewService(),
	}
}

func (h *Handler) ListLdap() iris.Handler {
	return func(ctx *context.Context) {
		ldaps, err := h.ldapService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", ldaps)
	}
}

func (h *Handler) AddLdap() iris.Handler {
	return func(ctx *context.Context) {
		var req Ldap
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.ldapService.Create(&req.Ldap, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) UpdateLdap() iris.Handler {
	return func(ctx *context.Context) {
		var req Ldap
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.ldapService.Update(req.UUID, &req.Ldap, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) TestConnect() iris.Handler {
	return func(ctx *context.Context) {
		var req Ldap
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		users, err := h.ldapService.TestConnect(&req.Ldap)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", users)
	}
}

func (h *Handler) SyncLdapUser() iris.Handler {
	return func(ctx *context.Context) {
		users, err := h.ldapService.GetLdapUser()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", users)
	}
}

func (h *Handler) TestLogin() iris.Handler {
	return func(ctx *context.Context) {
		var req TestLogin
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.ldapService.TestLogin(req.Username, req.Password)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) ImportUser() iris.Handler {
	return func(ctx *context.Context) {
		var req ImportRequest
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		result, err := h.ldapService.ImportUsers(req.Users)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", result)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/ldap")
	sp.Get("/", handler.ListLdap())
	sp.Post("/", handler.AddLdap())
	sp.Put("/", handler.UpdateLdap())
	sp.Post("/sync", handler.SyncLdapUser())
	sp.Post("/test/connect", handler.TestConnect())
	sp.Post("/test/login", handler.TestLogin())
	sp.Post("/import", handler.ImportUser())
}
