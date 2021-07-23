package cluster

import (
	"errors"
	"fmt"
	v1Project "github.com/KubeOperator/ekko/internal/model/v1/project"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func (h *Handler) CreateProject() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		var req v1Project.Project
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err)
			return
		}
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		k := kubernetes.NewKubernetes(*c)
		cert, err := k.CreateCommonUser(req.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		req.Certificate = cert
		if err := h.projectService.Create(c.Name, &req, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) ListProject() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", fmt.Sprintf("get cluster failed: %s", err.Error()))
			return
		}
		ps, err := h.projectService.List(c.Name, common.DBOptions{})
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		ctx.Values().Set("data", ps)
	}
}
