package session

import (
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func (h *Handler) UpdateProfile() iris.Handler {
	return func(ctx *context.Context) {
		var req ProfileSetter
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err)
			return
		}
		session := server.SessionMgr.Start(ctx)
		u := session.Get("profile")
		profile := u.(UserProfile)
		user, err := h.userService.GetByNameOrEmail(profile.Name, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		if req.Email != "" {
			user.Email = req.Email
		}
		if req.NickName != "" {
			user.NickName = req.NickName
		}
		if req.Language != "" {
			user.Language = req.Language
		}
		if err := h.userService.Update(profile.Name, user, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		profile = UserProfile{
			Name:                user.Name,
			NickName:            user.NickName,
			Email:               user.Email,
			Language:            user.Language,
			ResourcePermissions: profile.ResourcePermissions,
			IsAdministrator:     user.IsAdmin,
		}
		session.Set("profile", profile)
		ctx.Values().Set("data", "ok")
	}
}
func (h *Handler) UpdatePassword() iris.Handler {
	return func(ctx *context.Context) {
		var pass PasswordSetter
		if err := ctx.ReadJSON(&pass); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err)
			return
		}
		session := server.SessionMgr.Start(ctx)
		u := session.Get("profile")
		profile := u.(UserProfile)
		if err := h.userService.UpdatePassword(profile.Name, pass.OldPassword, pass.NewPassword, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", "can not match original password")
			return
		}
		ctx.Values().Set("data", "ok")
	}
}
