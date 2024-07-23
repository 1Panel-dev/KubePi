package sso

import (
	v1Session "github.com/KubeOperator/kubepi/internal/api/v1/session"
	v1Sso "github.com/KubeOperator/kubepi/internal/model/v1/sso"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/sso"
	"github.com/crewjam/saml/samlsp"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"net/http"
	"strings"
)

type Handler struct {
	ssoService sso.Service
}

var oauth2Config = &v1Sso.OpenID{}
var saml2Config = &samlsp.Middleware{}

func NewHandler() *Handler {
	return &Handler{
		ssoService: sso.NewService(),
	}
}

func (h *Handler) AddSso() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Sso.Sso
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.ssoService.Create(&req, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) ListSso() iris.Handler {
	return func(ctx *context.Context) {
		ssos, err := h.ssoService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", ssos)
	}
}

func (h *Handler) UpdateSso() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Sso.Sso
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		err := h.ssoService.Update(req.UUID, &req, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) LoginSso() iris.Handler {
	return func(ctx *context.Context) {
		ssos, err := h.ssoService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		// 根据协议设置重定向URL
		referer := ctx.GetHeader("Referer")
		callbakOpenID := strings.Replace(referer, "sso", "api/v1/sso/callback", -1)
		callbakSaml2 := strings.Replace(referer, "sso", "api/v1/sso/callback_saml2", -1)
		baseUrl := strings.Replace(referer, "sso", "api/v1/sso/", -1)

		// 目前只支持OpenID
		switch ssos.Protocol {
		case "openid":
			oauth2Config, err = h.ssoService.OpenIDConfig(ssos.ClientId, ssos.ClientSecret, ssos.InterfaceAddress, callbakOpenID)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			ctx.Redirect(oauth2Config.Oauth2Config.AuthCodeURL("state"), iris.StatusFound)
		case "saml2":
			saml2Config, err = h.ssoService.Saml2Config(ssos, baseUrl)
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			ctx.Redirect(callbakSaml2, iris.StatusFound)
		default:
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "目前只支持OpenID & SAML2 SSO认证~")
			return
		}
	}
}

func (h *Handler) CallbackOpenID() iris.Handler {
	return func(ctx *context.Context) {
		r := ctx.Request()
		language := ctx.GetHeader("Accept-Language")
		if strings.Contains(language, "zh-CN") {
			language = "zh-CN"
		} else {
			language = "en-US"
		}
		code := r.URL.Query().Get("code")

		oauth2Config.Code = code
		oauth2Config.Language = language
		userProfile, err := h.ssoService.OpenID(oauth2Config, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		// 默认为Session
		sId := ctx.GetCookie(server.SessionCookieName)
		if sId != "" {
			ctx.RemoveCookie(server.SessionCookieName)
			ctx.Request().Header.Del("Cookie")
		}
		sess := server.SessionMgr.Start(ctx)
		ctx.SetCookieKV(server.SessionCookieName, sess.ID())
		sess.Set("profile", userProfile)

		redirectURL := ""
		if strings.HasPrefix(strings.ToLower(r.Proto), "https") {
			redirectURL = "https://" + r.Host
		} else if strings.HasPrefix(strings.ToLower(r.Proto), "http") {
			redirectURL = "http://" + r.Host
		}
		ctx.Redirect(redirectURL, iris.StatusFound)
		handler := v1Session.NewHandler()
		go handler.SaveLoginLog(ctx, userProfile.Name)
		//ctx.Values().Set("data", userProfile)
	}
}

func (h *Handler) TestConnect() iris.Handler {
	return func(ctx *context.Context) {
		var req v1Sso.Sso
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
		}
		if err := h.ssoService.TestConnect(&req); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", "SSO测试连接成功")
	}
}

func (h *Handler) StatusSso() iris.Handler {
	return func(ctx *context.Context) {
		if ssoSwitch := h.ssoService.Status(common.DBOptions{}); !ssoSwitch {
			ctx.Values().Set("data", false)
			return
		}
		ctx.Values().Set("data", true)
	}
}

func (h *Handler) Saml2Auth() iris.Handler {
	return func(ctx *context.Context) {
		session, err := saml2Config.Session.GetSession(ctx.Request())
		if session != nil {
			r := ctx.Request()
			ctx.Request().WithContext(samlsp.ContextWithSession(ctx.Request().Context(), session))
			language := ctx.GetHeader("Accept-Language")
			if strings.Contains(language, "zh-CN") {
				language = "zh-CN"
			} else {
				language = "en-US"
			}

			attributes := session.(samlsp.SessionWithAttributes).GetAttributes()
			givenName := attributes.Get("givenName")
			email := attributes.Get("email")
			userProfile, err := h.ssoService.Saml2(givenName, email, language, common.DBOptions{})
			if err != nil {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			// 默认为Session
			sId := ctx.GetCookie(server.SessionCookieName)
			if sId != "" {
				ctx.RemoveCookie(server.SessionCookieName)
				ctx.Request().Header.Del("Cookie")
			}
			sess := server.SessionMgr.Start(ctx)
			ctx.SetCookieKV(server.SessionCookieName, sess.ID())
			sess.Set("profile", userProfile)

			redirectURL := ""
			if strings.HasPrefix(strings.ToLower(r.Proto), "https") {
				redirectURL = "https://" + r.Host
			} else if strings.HasPrefix(strings.ToLower(r.Proto), "http") {
				redirectURL = "http://" + r.Host
			}
			ctx.Redirect(redirectURL, iris.StatusFound)
			handler := v1Session.NewHandler()
			go handler.SaveLoginLog(ctx, userProfile.Name)
			//ctx.Values().Set("data", userProfile)
		}
		if err == samlsp.ErrNoSession {
			saml2Config.HandleStartAuthFlow(ctx.ResponseWriter(), ctx.Request())
			return
		}
		saml2Config.OnError(ctx.ResponseWriter(), ctx.Request(), err)
	}
}

func (h *Handler) Saml2Metadata() iris.Handler {
	return func(ctx *context.Context) {
		saml2Config.ServeMetadata(ctx.ResponseWriter(), ctx.Request())
	}
}

func (h *Handler) Saml2Acs() iris.Handler {
	return func(ctx *context.Context) {
		saml2Config.ServeACS(ctx.ResponseWriter(), ctx.Request())
	}
}

func (h *Handler) Saml2Slo() iris.Handler {
	return func(ctx *context.Context) {
		saml2Config.Session.DeleteSession(ctx.ResponseWriter(), ctx.Request())
		ctx.Redirect("/", http.StatusFound)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/sso")
	sp.Get("/", handler.ListSso())
	sp.Post("/", handler.AddSso())
	sp.Put("/", handler.UpdateSso())
	sp.Get("/login", handler.LoginSso())
	sp.Get("/callback", handler.CallbackOpenID())
	sp.Get("/callback_saml2", handler.Saml2Auth())
	sp.Post("/test/connect", handler.TestConnect())
	sp.Get("/status", handler.StatusSso())
	sp.Get("/saml/metadata", handler.Saml2Metadata())
	sp.Post("/saml/acs", handler.Saml2Acs())
	sp.Get("/saml/slo", handler.Saml2Slo())
}
