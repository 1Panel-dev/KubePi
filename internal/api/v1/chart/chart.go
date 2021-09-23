package chart

import (
	"errors"
	"github.com/KubeOperator/kubepi/internal/api/v1/session"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/chart"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	chartService chart.Service
}

func NewHandler() *Handler {
	return &Handler{
		chartService: chart.NewService(),
	}
}

func (h *Handler) SearchCharts() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		pattern := ctx.URLParam("pattern")
		charts, total, err := h.chartService.Search(pageNum, pageSize, pattern, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		cs := make([]Chart, 0)
		for i := range charts {
			cs = append(cs, Chart{
				Chart: charts[i],
			})
		}
		ctx.Values().Set("data", pkgV1.Page{Items: cs, Total: total})
	}
}

func (h *Handler) CreateChart() iris.Handler {
	return func(ctx *context.Context) {
		var req Chart
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
		if err := h.chartService.Create(&req.Chart, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteChart() iris.Handler {
	return func(ctx *context.Context) {
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}

		name := ctx.Params().GetString("name")
		if err := h.chartService.Delete(name, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
	}
}

func (h *Handler) GetChart() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.chartService.GetByName(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &Chart{Chart: *c})
	}
}
func (h *Handler) UpdateChart() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req Chart
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.chartService.Update(name, &req.Chart, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()
		ctx.Values().Set("data", &req)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/charts")
	sp.Post("/search", handler.SearchCharts())
	sp.Post("/", handler.CreateChart())
	sp.Delete("/:name", handler.DeleteChart())
	sp.Put("/:name", handler.UpdateChart())
	sp.Get("/:name", handler.GetChart())
}
