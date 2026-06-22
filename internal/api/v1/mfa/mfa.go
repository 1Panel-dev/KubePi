package mfa

import (
	"strings"

	sessionAuth "github.com/1Panel-dev/KubePi/internal/api/v1/session"
	"github.com/1Panel-dev/KubePi/internal/server"
	"github.com/1Panel-dev/KubePi/internal/service/v1/common"
	"github.com/1Panel-dev/KubePi/internal/service/v1/user"
	mfaUtil "github.com/1Panel-dev/KubePi/pkg/util/mfa"
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

func (m *Handler) MfaValidate() iris.Handler {
	return func(ctx *context.Context) {
		session := server.SessionMgr.Start(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(sessionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
			return
		}
		var mfa sessionAuth.MfaCredential
		if err := ctx.ReadJSON(&mfa); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u, err := m.userService.GetByNameOrEmail(p.Name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if strings.TrimSpace(u.Mfa.Secret) == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "mfa is not configured")
			return
		}
		success := mfaUtil.ValidCode(mfa.Code, u.Mfa.Secret)
		if !success {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "code is invalid")
			return
		} else {
			p.Mfa.Approved = true
			session.Set("profile", p)
			ctx.StatusCode(iris.StatusOK)
			return
		}
	}
}

func (m *Handler) MfaBind() iris.Handler {
	return func(ctx *context.Context) {
		session := server.SessionMgr.Start(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(sessionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
			return
		}
		var mfa sessionAuth.MfaCredential
		if err := ctx.ReadJSON(&mfa); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if strings.TrimSpace(mfa.Secret) == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "mfa secret is empty")
			return
		}
		success := mfaUtil.ValidCode(mfa.Code, mfa.Secret)
		if !success {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "code is invalid")
			return
		} else {
			us, err := m.userService.GetByNameOrEmail(p.Name, common.DBOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			if strings.TrimSpace(us.Mfa.Secret) != "" {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.Values().Set("message", "mfa is already configured")
				return
			}
			us.Mfa.Enable = true
			us.Mfa.Secret = mfa.Secret
			if err := m.userService.Update(p.Name, us, common.DBOptions{}); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			session.Delete("profile")
			ctx.StatusCode(iris.StatusOK)
			return
		}
	}
}

func (m *Handler) GetMfa() iris.Handler {
	return func(ctx *context.Context) {
		session := server.SessionMgr.Start(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(sessionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
			return
		}
		u, err := m.userService.GetByNameOrEmail(p.Name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if strings.TrimSpace(u.Mfa.Secret) != "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", "mfa is already configured")
			return
		}
		otp, err := mfaUtil.GetOtp(p.Name)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		} else {
			ctx.StatusCode(iris.StatusOK)
			ctx.Values().Set("data", otp)
			return
		}
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/mfa")
	sp.Get("/", handler.GetMfa())
	sp.Post("/bind", handler.MfaBind())
	sp.Post("/valid", handler.MfaValidate())
}
