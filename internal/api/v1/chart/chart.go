package chart

import (
	"bytes"
	"crypto/rand"
	"github.com/KubeOperator/kubepi/internal/service/v1/chart"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"math/big"
	"strings"
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
		//pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		//pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		//pattern := ctx.URLParam("pattern")
		//cluster := ctx.URLParam("cluster")
		//charts, total, err := h.chartService.Search(pageNum, pageSize, cluster, pattern, common.DBOptions{})
		//if err != nil {
		//	if !errors.Is(err, storm.ErrNotFound) {
		//		ctx.StatusCode(iris.StatusInternalServerError)
		//		ctx.Values().Set("message", err.Error())
		//		return
		//	}
		//}
		//cs := make([]Chart, 0)
		//for i := range charts {
		//	cs = append(cs, Chart{
		//		Chart: charts[i],
		//	})
		//}
		//ctx.Values().Set("data", pkgV1.Page{Items: cs, Total: total})
	}
}

//func (h *Handler) CreateChart() iris.Handler {
//	return func(ctx *context.Context) {
//		var req Chart
//		if err := ctx.ReadJSON(&req); err != nil {
//			ctx.StatusCode(iris.StatusBadRequest)
//			ctx.Values().Set("message", err.Error())
//			return
//		}
//		u := ctx.Values().Get("profile")
//		profile := u.(session.UserProfile)
//		//tx
//		tx, err := server.DB().Begin(true)
//		if err != nil {
//			ctx.StatusCode(iris.StatusInternalServerError)
//			ctx.Values().Set("message", err.Error())
//			return
//		}
//		req.CreatedBy = profile.Name
//		if err := h.chartService.Create(&req.Chart, common.DBOptions{DB: tx}); err != nil {
//			_ = tx.Rollback()
//			ctx.StatusCode(iris.StatusInternalServerError)
//			ctx.Values().Set("message", err.Error())
//			return
//		}
//		_ = tx.Commit()
//		ctx.Values().Set("data", req)
//	}
//}

func (h *Handler) DeleteChart() iris.Handler {
	return func(ctx *context.Context) {
		//tx, _ := server.DB().Begin(true)
		//txOptions := common.DBOptions{DB: tx}
		//
		//name := ctx.Params().GetString("name")
		//if err := h.chartService.Delete(name, txOptions); err != nil {
		//	_ = tx.Rollback()
		//	ctx.StatusCode(iris.StatusInternalServerError)
		//	ctx.Values().Set("message", err.Error())
		//	return
		//}
		//_ = tx.Commit()
	}
}

func (h *Handler) GetChart() iris.Handler {
	return func(ctx *context.Context) {
		//name := ctx.Params().GetString("name")
		//c, err := h.chartService.GetByName(name, common.DBOptions{})
		//if err != nil {
		//	ctx.StatusCode(iris.StatusInternalServerError)
		//	ctx.Values().Set("message", err.Error())
		//	return
		//}
		//ctx.Values().Set("data", &Chart{Chart: *c})
	}
}
func (h *Handler) UpdateChart() iris.Handler {
	return func(ctx *context.Context) {
	}
}

func (h *Handler) ListRepo() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.URLParam("cluster")
		entrys, err := h.chartService.SearchRepo(cluster)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		repos := make([]Repo, len(entrys))
		for i, en := range entrys {
			repos[i] = Repo{
				Name: en.Name,
				Url:  en.URL,
			}
		}
		ctx.Values().Set("data", repos)
	}
}

func (h *Handler) AddRepo() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.URLParam("cluster")
		var req RepoCreate
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}

		err := h.chartService.AddRepo(cluster, &req.RepoCreate)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) ListCharts() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		pattern := ctx.URLParam("pattern")
		cluster := ctx.URLParam("cluster")
		charts, total, err := h.chartService.ListCharts(cluster, pageNum, pageSize, pattern)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		chartArray := make([]*Chart, len(charts))
		for i, ch := range charts {
			names := strings.Split(ch.Name, "/")
			chartArray[i] = &Chart{
				Name:        ch.Chart.Metadata.Name,
				Description: ch.Chart.Metadata.Description,
				Icon:        ch.Chart.Icon,
				Repo:        names[0],
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: chartArray, Total: total})
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/charts")
	sp.Delete("/:name", handler.DeleteChart())
	sp.Put("/:name", handler.UpdateChart())
	sp.Get("/:name", handler.GetChart())
	sp.Get("/repos", handler.ListRepo())
	sp.Post("/repos", handler.AddRepo())
	sp.Post("/search", handler.ListCharts())
}
