package cluster

import (
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)


func (h *Handler) ListClusterRepos() iris.Handler {
	return func(ctx *context.Context) {
		cluster := ctx.Params().GetString("name")
		clusterRepos, err := h.clusterRepoService.List(cluster, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound  {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", clusterRepos)
	}
}

func (h *Handler) AddCLusterRepo() iris.Handler {
	return func(ctx *context.Context) {
		var req ClusterRepo
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.clusterRepoService.Create(&req.ClusterRepo, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}
//
//func Install(parent iris.Party) {
//	handler := NewHandler()
//	sp := parent.Party("/repos/:cluster")
//	sp.Get("/", handler.ListClusterRepo())
//	sp.Post("/", handler.AddCLusterRepo())
//}
