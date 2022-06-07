package file

import (
	"bufio"
	"fmt"
	fileModel "github.com/KubeOperator/kubepi/internal/model/v1/file"
	"github.com/KubeOperator/kubepi/internal/service/v1/file"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
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

func (h *Handler) CreateFolder() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.Commands = []string{"mkdir", req.Path}
		if _, err := h.fileService.ExecNewCommand(req); err != nil {
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
		command := "echo '" + req.Content + "' >> " + req.Path
		req.Commands = []string{"sh", "-c", command}
		if _, err := h.fileService.ExecNewCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) UpdateFile() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.fileService.EditFile(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) OpenFile() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.Commands = []string{"cat", req.Path}
		res, err := h.fileService.ExecNewCommand(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", string(res))
	}
}

func (h *Handler) ReNameFile() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if req.OldPath == "" {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "file or path is not exist")
			return
		}
		req.Commands = []string{"mv", req.OldPath, req.Path}
		_, err := h.fileService.ExecNewCommand(req)
		if err != nil {
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
		req.Commands = []string{"rm", req.Path}
		if _, err := h.fileService.ExecNewCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) DownloadFile() iris.Handler {
	return func(ctx *context.Context) {

		var req fileModel.Request
		req.Path = ctx.URLParam("path")
		req.Namespace = ctx.URLParam("namespace")
		req.Cluster = ctx.URLParam("cluster")
		req.PodName = ctx.URLParam("podName")
		req.ContainerName = ctx.URLParam("containerName")

		file, err := h.fileService.DownloadFile(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		filename := path.Base(file)
		err = ctx.SendFile(file, filename)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		os.RemoveAll(file)
	}
}

func (h *Handler) UploadFile() iris.Handler {
	return func(ctx *context.Context) {
		f, header, err := ctx.FormFile("file")
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		var req fileModel.Request
		req.Path = ctx.URLParam("path")
		req.Namespace = ctx.URLParam("namespace")
		req.Cluster = ctx.URLParam("cluster")
		req.PodName = ctx.URLParam("podName")
		req.ContainerName = ctx.URLParam("containerName")

		path := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		path = filepath.Join(path, "/"+header.Filename)
		file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		r := bufio.NewReader(f)
		w := bufio.NewWriter(file)

		size := 4 * 1024
		buf := make([]byte, 4*1024)
		for {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			if n == 0 {
				break
			}
			_, err = w.Write(buf[:n])
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			if n < size {
				break
			}
		}
		err = w.Flush()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		req.FilePath = path
		if req.Path == "/" {
			req.Path = req.Path + header.Filename
		} else {
			req.Path = req.Path + "/" + header.Filename
		}
		err = h.fileService.UploadFile(req)
		if err != nil {
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
	sp.Post("/files/open", handler.OpenFile())
	sp.Post("/files/rename", handler.ReNameFile())
	sp.Post("/files/upload", handler.UploadFile())
	sp.Post("/files/update", handler.UpdateFile())
	sp.Get("/files/download", handler.DownloadFile())
}
