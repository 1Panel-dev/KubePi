package imagerepo

import (
	"errors"

	"github.com/KubeOperator/kubepi/internal/api/v1/commons"
	v1ImageRepo "github.com/KubeOperator/kubepi/internal/model/v1/imagerepo"
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

// Create Repo
// @Tags repos
// @Summary Create repo
// @Description Create repo
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Success 200 {object} RepoConfig
// @Security ApiKeyAuth
// @Router /imagerepos [post]
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

// List Internal Repos
// @Tags repos
// @Summary List Internal repos
// @Description List Internal repos
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Success 200 {object} []string
// @Security ApiKeyAuth
// @Router /imagerepos/repositories/search [post]
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

// Update Repo
// @Tags repos
// @Summary Update repo by name
// @Description Update repo by name
// @Accept  json
// @Produce  json
// @Param request body RepoConfig true "request"
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [put]
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

// Delete Repo
// @Tags repos
// @Summary Delete repo by name
// @Description Delete repo by name
// @Accept  json
// @Produce  json
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [delete]
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

// Get Repo
// @Tags repos
// @Summary Get repo by name
// @Description Get repo by name
// @Accept  json
// @Produce  json
// @Param name path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/{name} [get]
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

// List Repo for Cluster
// @Tags repos
// @Summary List repo for cluster
// @Description Get repo for cluster
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Success 200 {object} []v1ImageRepo.ImageRepo
// @Security ApiKeyAuth
// @Router /imagerepos/cluster/{cluster} [get]
func (h *Handler) ListRepoForCluster() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("cluster")
		var imageRepos []v1ImageRepo.ImageRepo
		imageRepos, err := h.imageRepoService.ListByCluster(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", imageRepos)
	}
}

// List Images for Cluster
// @Tags repos
// @Summary List images for cluster
// @Description Get images for cluster
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Param repo path string true "镜像仓库名称"
// @Success 200 {object} []string
// @Security ApiKeyAuth
// @Router /imagerepos/{cluster}/{repo} [put]
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

// List Images By Repo
// @Tags repos
// @Summary List images by repo
// @Description Get images by repo
// @Accept  json
// @Produce  json
// @Param repo path string true "镜像仓库名称"
// @Success 200 {object} v1ImageRepo.RepoResponse
// @Security ApiKeyAuth
// @Router /imagerepos/{repo}/search [get]
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
