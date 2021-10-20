package chart

import (
	"github.com/KubeOperator/kubepi/internal/service/v1/chart"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
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

func (h *Handler) DeleteRepo() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.URLParam("cluster")
		name := ctx.Params().GetString("name")
		if err := h.chartService.RemoveRepo(cluster, name); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) GetChart() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		cluster := ctx.URLParam("cluster")
		repo := ctx.URLParam("repo")
		cs, err := h.chartService.GetCharts(cluster, repo, name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", cs)
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
		repo := ctx.URLParam("repo")
		charts, total, err := h.chartService.ListCharts(cluster, repo, pageNum, pageSize, pattern)
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

func (h *Handler) GetChartByVersion() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		cluster := ctx.URLParam("cluster")
		repo := ctx.URLParam("repo")
		version := ctx.URLParam("version")
		cs, err := h.chartService.GetChartByVersion(cluster, repo, name, version)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", cs)
	}
}

func (h *Handler) GetChartForUpdate() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		cluster := ctx.URLParam("cluster")
		repo := ctx.URLParam("repo")
		cs, err := h.chartService.GetChartsUpdate(cluster, repo, name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", cs)
	}
}

func (h *Handler) InstallChart() iris.Handler {
	return func(ctx *context.Context) {
		var req ChInstall
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.chartService.InstallChart(req.Cluster, req.Repo, req.Namespace, req.Name, req.ChartName, req.ChartVersion, req.Values)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) UpdateChart() iris.Handler {
	return func(ctx *context.Context) {
		var req ChInstall
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.chartService.UpgradeChart(req.Cluster, req.Repo, req.Name, req.ChartName, req.ChartVersion, req.Values)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) AllInstalled() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		pattern := ctx.URLParam("pattern")
		cluster := ctx.URLParam("cluster")
		installed, total, err := h.chartService.ListAllInstalled(cluster, pageNum, pageSize, pattern)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", pkgV1.Page{Items: installed, Total: total})
	}
}

func (h *Handler) UnInstall() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.URLParam("cluster")
		name := ctx.Params().GetString("name")
		err := h.chartService.UnInstallChart(cluster, name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", "")
	}
}

func (h *Handler) GetAppDetail() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.URLParam("cluster")
		name := ctx.Params().GetString("name")
		data, err := h.chartService.GetAppDetail(cluster, name)
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
	sp := parent.Party("/charts")
	sp.Put("/:name", handler.UpdateChart())
	sp.Get("/:name", handler.GetChart())
	sp.Get("/repos", handler.ListRepo())
	sp.Post("/repos", handler.AddRepo())
	sp.Delete("/repos/:name", handler.DeleteRepo())
	sp.Post("/search", handler.ListCharts())
	sp.Get("/detail/:name", handler.GetChartByVersion())
	sp.Post("/install", handler.InstallChart())
	app := parent.Party("/apps")
	app.Post("/search", handler.AllInstalled())
	app.Delete("/:name", handler.UnInstall())
	app.Get("/:name", handler.GetAppDetail())
	app.Get("/update/:name", handler.GetChartForUpdate())
	app.Put("/upgrade/:name", handler.UpdateChart())
}
