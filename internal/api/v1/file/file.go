package file

import (
	fileModel "github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/file"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"strings"
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
		req.Commands = []string{"./kotools", "ls", req.Path}
		res, err := h.fileService.ListFiles(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", res)
	}
}

func (h *Handler) CreateFolder() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.Commands = []string{"./kotools", "mkdir", req.Path}
		if err := h.fileService.ExecCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) CreateFile() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.Commands = []string{"./kotools", "touch", req.Path}
		req.Stdin = strings.NewReader(req.Content)
		if err := h.fileService.ExecCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) RmFolder() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.Commands = []string{"./kotools", "rm", req.Path}
		if err := h.fileService.ExecCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/pod")
	sp.Post("/files", handler.ListFiles())
	sp.Post("/folder/create", handler.CreateFolder())
	sp.Post("/folder/delete", handler.RmFolder())
	sp.Post("/files/create", handler.CreateFile())
}
