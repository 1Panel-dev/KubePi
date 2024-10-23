package monitor

import (
	"errors"
	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	v1Monitor "github.com/KubeOperator/kubepi/internal/model/v1/monitor"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/monitor"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/asdine/storm/v3"
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

func (h *Handler) SearchMetrics() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		metrics, total, err := h.monitorService.MetricsSearch(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: metrics, Total: total})
	}
}

func (h *Handler) AddMetrics() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.MetricsConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.monitorService.MetricsCreate(&req, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteMetrics() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}
		if err = h.monitorService.MetricsDelete(name, txOptions); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}
}

func (h *Handler) GetMetrics() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		metrics, err := h.monitorService.MetricsGetByName(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", metrics)
	}
}

func (h *Handler) UpdateMetrics() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Monitor.MetricsConfig
		name := ctx.Params().GetString("name")
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.monitorService.MetricsUpdate(name, &req, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) ExplorerMetrics() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		metrics, err := h.monitorService.MetricsExplorer(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", metrics)
	}
}

func (h *Handler) TestConnectMetrics() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		if err := h.monitorService.MetricsTestConnect(name, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", name+"连通性测试成功")
	}
}

func (h *Handler) QueryMetrics() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		promql := ctx.URLParam("promql")
		timestamp := ctx.URLParam("time")

		data, err := h.monitorService.MetricsQuery(name, promql, timestamp, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", data)
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

	metrics := parent.Party("/monitor/metrics")
	metrics.Post("/", handler.AddMetrics())
	metrics.Delete("/:name", handler.DeleteMetrics())
	metrics.Post("/search", handler.SearchMetrics())
	metrics.Get("/:name", handler.GetMetrics())
	metrics.Put("/:name", handler.UpdateMetrics())
	metrics.Get("/:name/explorer", handler.ExplorerMetrics())
	metrics.Get("/:name/test/connect", handler.TestConnectMetrics())
}
