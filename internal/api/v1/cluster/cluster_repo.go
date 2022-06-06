package cluster

import (
	"github.com/KubeOperator/kubepi/internal/model/v1/clusterrepo"
	_ "github.com/KubeOperator/kubepi/internal/model/v1/imagerepo"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// List ClusterRepos
// @Tags clusters
// @Summary List all clusterRepos
// @Description List all clusterRepos
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Success 200 {object} []clusterrepo.ClusterRepo
// @Security ApiKeyAuth
// @Router /clusters/{cluster}/repos [get]
func (h *Handler) ListClusterRepos() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("name")
		clusterRepos, err := h.clusterRepoService.List(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", clusterRepos)
	}
}

// Get ClusterRepo Detail
// @Tags clusters
// @Summary Get ClusterRepo Detail
// @Description Get ClusterRepo Detail
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Success 200 {object} []imagerepo.ImageRepo
// @Security ApiKeyAuth
// @Router /clusters/{cluster}/repos [get]
func (h *Handler) ListClusterReposDetail() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("name")
		clusterRepos, err := h.clusterRepoService.ListInfo(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", clusterRepos)
	}
}

// Create Cluster Repo
// @Tags clusters
// @Summary Create Cluster Repo
// @Description Create Cluster Repo
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Param request body CreateRepo true "request"
// @Success 200 {object} CreateRepo
// @Security ApiKeyAuth
// @Router /clusters/{cluster}/repos [post]
func (h *Handler) AddCLusterRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req CreateRepo
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
		txOptions := common.DBOptions{DB: tx}
		for _, v := range req.Repos {
			clusterRepo := &clusterrepo.ClusterRepo{
				Cluster: req.Cluster,
				Repo:    v,
			}
			err := h.clusterRepoService.Create(clusterRepo, txOptions)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				_ = tx.Rollback()
				return
			}
		}
		_ = tx.Commit()
		ctx.Values().Set("data", &req)
	}
}

// Delete ClusterRepo
// @Tags clusters
// @Summary Delete clusterRepo by name
// @Description Delete clusterRepo by name
// @Accept  json
// @Produce  json
// @Param cluster path string true "集群名称"
// @Param repo path string true "镜像仓库名称"
// @Success 200 {number} 200
// @Security ApiKeyAuth
// @Router /clusters/{cluster}/repos/{repo} [delete]
func (h *Handler) DeleteClusterRepo() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("name")
		repo := ctx.Params().GetString("repo")
		err := h.clusterRepoService.Delete(cluster, repo, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", "")
	}
}
