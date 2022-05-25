package imagerepo

import (
	"errors"

	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterrepo"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/imagerepo"
	pkgV1 "github.com/KubeOperator/kubepi/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	imageRepoService   imagerepo.Service
	clusterRepoService clusterrepo.Service
}

func NewHandler() *Handler {
	return &Handler{
		imageRepoService:   imagerepo.NewService(),
		clusterRepoService: clusterrepo.NewService(),
	}
}

func (h *Handler) SearchRepos() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions commons.SearchConditions
		if err := ctx.ReadJSON(&conditions); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		repos, total, err := h.imageRepoService.Search(pageNum, pageSize, conditions.Conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: repos, Total: total})
	}
}
func (h *Handler) CreateRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.imageRepoService.Create(&req.ImageRepo, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) ListInternalRepos() iris.Handler {
	return func(ctx *context.Context) {
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		names, err := h.imageRepoService.ListInternalRepos(req.ImageRepo, req.Page, req.Limit, req.Search)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", names)
	}
}

func (h *Handler) UpdateRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req RepoConfig
		imageRepoName := ctx.Params().GetString("name")
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.imageRepoService.UpdateRepo(imageRepoName, &req.ImageRepo, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) DeleteRepo() iris.Handler {
	return func(ctx *context.Context) {
		imageRepoName := ctx.Params().GetString("name")
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}
		if err := h.imageRepoService.Delete(imageRepoName, txOptions); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			_ = tx.Rollback()
			return
		}
		if err := h.clusterRepoService.DeleteByRepo(imageRepoName, txOptions); err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}
}

func (h *Handler) GetRepo() iris.Handler {
	return func(ctx *context.Context) {
		imageRepoName := ctx.Params().GetString("name")
		imageRepo, err := h.imageRepoService.GetByName(imageRepoName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", imageRepo)
	}
}

func (h *Handler) ListRepoForCluster() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("cluster")
		imageRepos, err := h.imageRepoService.ListByCluster(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", imageRepos)
	}
}

func (h *Handler) ListImages() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("cluster")
		name := ctx.Params().GetString("repo")
		imageRepos, err := h.imageRepoService.ListImages(name, cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", imageRepos)
	}
}

func (h *Handler) ListImagesByRepo() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("repo")
		var req RepoConfig
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		res, err := h.imageRepoService.ListImagesByRepo(name, req.Page, req.Limit, req.Search, req.ContinueToken, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", res)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/imagerepos")
	sp.Post("/search", handler.SearchRepos())
	sp.Post("/", handler.CreateRepo())
	sp.Delete("/:name", handler.DeleteRepo())
	sp.Post("/repositories/search", handler.ListInternalRepos())
	sp.Get("/:name", handler.GetRepo())
	sp.Put("/:name", handler.UpdateRepo())
	sp.Get("/cluster/:cluster", handler.ListRepoForCluster())
	sp.Get("/images/:cluster/:repo", handler.ListImages())
	sp.Post("/images/:repo/search", handler.ListImagesByRepo())
}
