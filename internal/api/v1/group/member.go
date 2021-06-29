package group

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func (h *Handler) GetGroupMembers() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		bindings, err := h.groupBindingService.ListByGroupName(groupName, common.DBOptions{})
		if err != nil {
			if !errors.As(err, &storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		var members []Member

		for i := range bindings {
			members = append(members, Member{
				Name:     bindings[i].UserRef,
				Group:    bindings[i].GroupRef,
				CreateAt: bindings[i].CreateAt,
				CreateBy: bindings[i].CreatedBy,
			})
		}

		ctx.Values().Set("data", members)
	}
}

func (h *Handler) CreateGroupMember() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		var req Member
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		binding := v1Group.Binding{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "GroupBinding",
				CreatedBy:  profile.Name,
			},
			Metadata: v1.Metadata{
				Name: fmt.Sprintf("%s-%s-binding", req.Name, groupName),
			},
			UserRef:  req.Name,
			GroupRef: groupName,
		}
		if err := h.groupBindingService.Create(&binding, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) DeleteGroupMember() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		memberName := ctx.Params().GetString("member")
		gb, err := h.groupBindingService.GetGroupBindingByGroupNameAndUserName(memberName, groupName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.groupBindingService.Delete(gb.Name, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

	}
}
