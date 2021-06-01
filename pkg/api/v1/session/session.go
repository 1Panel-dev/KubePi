package session

import (
	"github.com/KubeOperator/ekko/pkg/service/v1/user"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	userService user.Service
}

func NewHandler() *Handler {
	return &Handler{userService: user.NewService()}
}

func (h *Handler) Login() iris.Handler {
	return func(ctx *context.Context) {

	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("sessions")
	sp.Post("login", handler.Login())
}
