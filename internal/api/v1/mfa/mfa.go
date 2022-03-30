package mfa

import (
	serrionAuth "github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1User "github.com/KubeOperator/kubepi/internal/model/v1/user"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/user"
	mfaUtil "github.com/KubeOperator/kubepi/pkg/util/mfa"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
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
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(serrionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
			return
		}
		var mfa serrionAuth.MfaCredential
		if err := ctx.ReadJSON(&mfa); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		success := mfaUtil.ValidCode(mfa.Code, mfa.Secret)
		if !success {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "code is not valid")
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
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(serrionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
			return
		}
		var mfa serrionAuth.MfaCredential
		if err := ctx.ReadJSON(&mfa); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		success := mfaUtil.ValidCode(mfa.Code, mfa.Secret)
		if !success {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "code is not valid")
			return
		} else {
			session.Delete("profile")
			us := &v1User.User{
				Metadata: v1.Metadata{
					Name: mfa.Username,
				},
				Mfa: v1User.Mfa{
					Enable: true,
					Secret: mfa.Secret,
				},
			}
			if err := m.userService.Update(mfa.Username, us, common.DBOptions{}); err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			ctx.StatusCode(iris.StatusOK)
			return
		}
	}
}

func (m *Handler) GetMfa() iris.Handler {
	return func(ctx *context.Context) {
		session := sessions.Get(ctx)
		loginUser := session.Get("profile")
		if loginUser == nil {
			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.Values().Set("message", "no login user")
			return
		}
		p, ok := loginUser.(serrionAuth.UserProfile)
		if !ok {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not parse to session user")
			return
		}
		if p.Mfa.Enable == false {
			ctx.StatusCode(iris.StatusOK)
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
