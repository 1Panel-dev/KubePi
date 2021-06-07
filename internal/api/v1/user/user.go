package user

import (
	"errors"
	"github.com/KubeOperator/ekko/internal/service/v1/user"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	userService user.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService: user.NewService(),
	}
}

func (h *Handler) SearchUsers() iris.Handler {
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
		users, total, err := h.userService.Search(pageNum, pageSize, conditions)
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: users, Total: total})
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/users")
	sp.Post("/search", handler.SearchUsers())
}
