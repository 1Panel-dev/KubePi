package session

import (
	"fmt"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/internal/service/v1/user"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	userService user.Service
}

func NewHandler() *Handler {
	return &Handler{
		userService: user.NewService(),
	}
}

func (h *Handler) IsLogin() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", loginUser != nil)
	}
}

func (h *Handler) Login() iris.Handler {
	return func(ctx *context.Context) {
		var loginCredential LoginCredential
		if err := ctx.ReadJSON(&loginCredential); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u, err := h.userService.Get(loginCredential.Username, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("query user %s error: %s", loginCredential.Username, err.Error()))
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.Spec.Authenticate.Password), []byte(loginCredential.Password)); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "username or password error")
			return
		}
		session := sessions.Get(ctx)
		profile := UserProfile{
			Name:     u.Name,
			NickName: u.Spec.Info.NickName,
			Email:    u.Spec.Info.Email,
		}
		session.Set("profile", profile)
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", profile)
	}
}

func (h *Handler) Logout() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		session.Delete("profile")
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", "logout success")
	}
}

func (h *Handler) GetProfile() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", loginUser)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/sessions")
	sp.Post("", handler.Login())
	sp.Delete("", handler.Logout())
	sp.Get("", handler.GetProfile())
	sp.Get("/status", handler.IsLogin())
}
