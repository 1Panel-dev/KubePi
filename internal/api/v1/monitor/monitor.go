package monitor

import (
	v1Monitor "github.com/KubeOperator/kubepi/internal/model/v1/monitor"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/monitor"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	monitorService monitor.Service
}

func NewHandler() *Handler {
	return &Handler{
		monitorService: monitor.NewService(),
	}
}

func (h *Handler) AddGrafana() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.GrafanaConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.monitorService.GrafanaCreate(&req, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) ListGrafana() iris.Handler {
	return func(ctx *context.Context) {
		monitor, err := h.monitorService.GrafanaList(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", monitor)
	}
}

func (h *Handler) UpdateGrafana() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.GrafanaConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.monitorService.GrafanaUpdate(req.UUID, &req, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) TestConnectGrafana() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.GrafanaConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		if err := h.monitorService.GrafanaTestConnect(&req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", "Grafana连通性测试成功")
	}
}

func (h *Handler) ImportDashboardsGrafana() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.GrafanaConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.monitorService.GrafanaImportDashboards(&req)
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
	grafana := parent.Party("/monitor/grafana")
	grafana.Get("/", handler.ListGrafana())
	grafana.Post("/", handler.AddGrafana())
	grafana.Put("/", handler.UpdateGrafana())
	grafana.Post("/test/connect", handler.TestConnectGrafana())
	grafana.Post("/import", handler.ImportDashboardsGrafana())
}
