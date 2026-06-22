package file

import (
	"archive/tar"
	"errors"
	"fmt"
	"github.com/1Panel-dev/KubePi/internal/api/v1/session"
	fileModel "github.com/1Panel-dev/KubePi/internal/model/v1/file"
	"github.com/1Panel-dev/KubePi/internal/service/v1/file"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
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

func withRequestUser(ctx *context.Context, req *fileModel.Request) {
	u := ctx.Values().Get("profile")
	profile := u.(session.UserProfile)
	req.UserName = profile.Name
	req.IsAdmin = profile.IsAdministrator
}

func (h *Handler) ListFiles() iris.Handler {
	return func(ctx *context.Context) {
		var req fileModel.Request
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		withRequestUser(ctx, &req)
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
		withRequestUser(ctx, &req)
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
		withRequestUser(ctx, &req)
		req.Stdin = strings.NewReader(req.Content)
		req.Commands = []string{"tee", req.Path}
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
		withRequestUser(ctx, &req)
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
		withRequestUser(ctx, &req)
		res, err := h.fileService.CatFile(req)
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
		withRequestUser(ctx, &req)
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
		withRequestUser(ctx, &req)
		req.Commands = []string{"rm", req.Path}
		if _, err := h.fileService.ExecNewCommand(req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func (h *Handler) DownloadFolder() iris.Handler {
	return func(ctx *context.Context) {

		var req fileModel.Request
		req.Path = ctx.URLParam("path")
		req.Namespace = ctx.URLParam("namespace")
		req.Cluster = ctx.URLParam("cluster")
		req.PodName = ctx.URLParam("podName")
		req.ContainerName = ctx.URLParam("containerName")
		withRequestUser(ctx, &req)

		file, err := h.fileService.DownloadFolder(req)
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
func (h *Handler) DownloadFile() iris.Handler {
	return func(ctx *context.Context) {

		var req fileModel.Request
		req.Path = ctx.URLParam("path")
		req.Namespace = ctx.URLParam("namespace")
		req.Cluster = ctx.URLParam("cluster")
		req.PodName = ctx.URLParam("podName")
		req.ContainerName = ctx.URLParam("containerName")
		withRequestUser(ctx, &req)

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
		var req fileModel.Request
		req.Path = ctx.URLParam("path")
		req.Namespace = ctx.URLParam("namespace")
		req.Cluster = ctx.URLParam("cluster")
		req.PodName = ctx.URLParam("podName")
		req.ContainerName = ctx.URLParam("containerName")
		withRequestUser(ctx, &req)

		srcPath := filepath.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixNano()))
		err := saveTarFile(ctx, srcPath)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		req.FilePath = srcPath
		err = h.fileService.UploadFile(req)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
	}
}

func saveTarFile(ctx *context.Context, srcPath string) error {
	maxSize := ctx.Application().ConfigurationReadOnly().GetPostMaxMemory()
	err := ctx.Request().ParseMultipartForm(maxSize)
	if err != nil {
		return err
	}
	form := ctx.Request().MultipartForm
	files := form.File["files"]
	if len(files) == 0 {
		return errors.New("files is null")
	}
	fw, err := os.Create(srcPath)
	defer fw.Close()
	if err != nil {
		return err
	}
	tw := tar.NewWriter(fw)
	defer tw.Close()
	for _, f := range files {
		hdr := &tar.Header{
			Name: f.Filename,
			Mode: 0644,
			Size: f.Size,
		}
		err = tw.WriteHeader(hdr)
		if err != nil {
			return err
		}
		_f, err := f.Open()
		if err != nil {
			return err
		}
		_, err = io.Copy(tw, _f)
		if err != nil {
			return err
		}
	}
	return nil
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
	sp.Get("/files/download/folder", handler.DownloadFolder())
	sp.Get("/files/download/file", handler.DownloadFile())
}
