package cluster

import (
	"github.com/KubeOperator/kubepi/internal/model/v1/clusterrepo"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

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
				tx.Rollback()
				return
			}
		}
		tx.Commit()
		ctx.Values().Set("data", &req)
	}
}

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
