package group

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/clusterbinding"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/internal/service/v1/group"
	"github.com/KubeOperator/ekko/internal/service/v1/groupbinding"
	"github.com/KubeOperator/ekko/internal/service/v1/rolebinding"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/KubeOperator/ekko/pkg/collectons"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	groupService          group.Service
	groupBindingService   groupbinding.Service
	roleBindingService    rolebinding.Service
	clusterBindingService clusterbinding.Service
	clusterService        cluster.Service
}

func NewHandler() *Handler {
	return &Handler{
		groupService:          group.NewService(),
		groupBindingService:   groupbinding.NewService(),
		roleBindingService:    rolebinding.NewService(),
		clusterBindingService: clusterbinding.NewService(),
		clusterService:        cluster.NewService(),
	}
}

func (h *Handler) SearchGroups() iris.Handler {
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
		groups, total, err := h.groupService.Search(pageNum, pageSize, conditions, common.DBOptions{})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		ctx.Values().Set("data", pkgV1.Page{Items: groups, Total: total})
	}
}

func (h *Handler) CreateGroup() iris.Handler {
	return func(ctx *context.Context) {
		var req Group
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.CreatedBy = profile.Name
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		txOptions := common.DBOptions{DB: tx}
		if err := h.groupService.Create(&req.Group, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		for i := range req.Roles {
			roleName := req.Roles[i]
			binding := v1Role.Binding{
				BaseModel: v1.BaseModel{
					ApiVersion: "v1",
					Kind:       "RoleBinding",
				},
				Metadata: v1.Metadata{
					Name: fmt.Sprintf("%s-group:%s-binding", roleName, req.Name),
				},
				Subject: v1Role.Subject{
					Kind: "Group",
					Name: req.Name,
				},
				RoleRef: roleName,
			}
			if err := h.roleBindingService.CreateRoleBinding(&binding, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		_ = tx.Commit()
		ctx.Values().Set("data", req)
	}
}

func (h *Handler) UpdateGroup() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		var req Group
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}

		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		tx, err := server.DB().Begin(true)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		if err := h.groupService.Update(groupName, &req.Group, common.DBOptions{DB: tx}); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		bindings, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{Kind: "Group", Name: groupName}, common.DBOptions{})
		if err != nil {
			if !errors.As(err, &storm.ErrNotFound) {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		currentRoles := collectons.NewStringSet()
		for i := range bindings {
			currentRoles.Add(bindings[i].RoleRef)
		}
		for i := range req.Roles {
			r := req.Roles[i]
			if currentRoles.Exists(r) {
				continue
			}
			binding := v1Role.Binding{
				BaseModel: v1.BaseModel{
					Kind:       "RoleBind",
					ApiVersion: "v1",
					CreatedBy:  profile.Name,
				},
				Metadata: v1.Metadata{
					Name: fmt.Sprintf("role-binding-%s-%s", r, req.Name),
				},
				Subject: v1Role.Subject{
					Kind: "Group",
					Name: req.Name,
				},
				RoleRef: r,
			}
			if err := h.roleBindingService.CreateRoleBinding(&binding, common.DBOptions{DB: tx}); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
			currentRoles.Add(binding.RoleRef)
		}
		diffs := currentRoles.Difference(req.Roles)

		for i := range bindings {
			for j := range diffs {
				if bindings[i].RoleRef == diffs[j] {
					if err := h.roleBindingService.Delete(bindings[i].Name, common.DBOptions{DB: tx}); err != nil {
						_ = tx.Rollback()
						ctx.StatusCode(iris.StatusInternalServerError)
						ctx.Values().Set("message", err.Error())
						return
					}
				}
			}
		}
		_ = tx.Commit()
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) DeleteGroup() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		tx, _ := server.DB().Begin(true)
		txOptions := common.DBOptions{DB: tx}

		gbs, err := h.groupBindingService.ListByGroupName(groupName, txOptions)
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}

		var affectUserNames []string

		for i := range gbs {
			affectUserNames = append(affectUserNames, gbs[i].UserRef)
			if err := h.groupBindingService.Delete(gbs[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}

		cbs, err := h.clusterBindingService.GetBindingsBySubject(v1Cluster.Subject{
			Kind: "Group",
			Name: groupName,
		}, txOptions)
		if err != nil && !errors.As(err, &storm.ErrNotFound) {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		for i := range cbs {
			if err := h.clusterBindingService.Delete(cbs[i].Name, txOptions); err != nil {
				_ = tx.Rollback()
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		if err := h.groupService.Delete(groupName, txOptions); err != nil {
			_ = tx.Rollback()
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		_ = tx.Commit()

		go func() {
			for i := range affectUserNames {
				h.modifyClusterUserCert(affectUserNames[i])
			}
		}()
	}
}

func (h *Handler) GetGroup() iris.Handler {
	return func(ctx *context.Context) {
		groupName := ctx.Params().GetString("name")
		if groupName == "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("invalid resource name %s", groupName))
			return
		}
		g, err := h.groupService.Get(groupName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		bindings, err := h.roleBindingService.GetRoleBindingBySubject(v1Role.Subject{Kind: "Group", Name: groupName}, common.DBOptions{})
		gp := Group{
			Group: *g,
			Roles: []string{},
		}
		for i := range bindings {
			gp.Roles = append(gp.Roles, bindings[i].RoleRef)
		}
		ctx.Values().Set("data", &gp)
	}
}

func (h *Handler) GetGroups() iris.Handler {
	return func(ctx *context.Context) {
		us, err := h.groupService.List(common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", us)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/groups")
	sp.Post("/search", handler.SearchGroups())
	sp.Get("/", handler.GetGroups())
	sp.Post("/", handler.CreateGroup())
	sp.Put("/:name", handler.UpdateGroup())
	sp.Get("/:name", handler.GetGroup())
	sp.Delete("/:name", handler.DeleteGroup())
	sp.Get("/:name/members", handler.GetGroupMembers())
	sp.Post("/:name/members", handler.CreateGroupMember())
	sp.Delete("/:name/members/:member", handler.DeleteGroupMember())

}
