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

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/ldap")
	sp.Get("/", handler.ListLdap())
	sp.Post("/", handler.AddLdap())
	sp.Put("/", handler.UpdateLdap())
}
