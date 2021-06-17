package group

import (
	"errors"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/internal/service/v1/group"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	groupService group.Service
}

func NewHandler() *Handler {
	return &Handler{
		groupService: group.NewService(),
	}
}

func (h *Handler) SearchGroups() iris.Handler {
	return func(ctx *context.Context) {
		pageNum, _ := ctx.Values().GetInt(pkgV1.PageNum)
		pageSize, _ := ctx.Values().GetInt(pkgV1.PageSize)
		var conditions pkgV1.Conditions
		if ctx.GetContentLength() > 0 {
			if err := ctx.ReadJSON(&conditions); err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		groups, total, err := h.groupService.Search(pageNum, pageSize, conditions)
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: groups, Total: total})
	}
}

func (h *Handler) GetGroups() iris.Handler {
	return func(ctx *context.Context) {
		us, err := h.groupService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", us)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/groups")
	sp.Post("/search", handler.SearchGroups())
	sp.Get("/", handler.GetGroups())
}
