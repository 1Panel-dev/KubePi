package file

import (
	fileModel "github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/file"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	fileService file.Service
}

func NewHandler() *Handler {
	return &Handler{
		fileService: file.NewService(),
	}
}

func (h *Handler) ListFiles() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		res, err := h.fileService.ListFiles(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", res)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/pod/files")
	sp.Post("", handler.ListFiles())
}
