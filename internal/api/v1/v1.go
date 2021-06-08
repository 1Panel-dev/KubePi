package v1

import (
	"errors"
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/cluster"
	"github.com/KubeOperator/ekko/internal/api/v1/group"
	"github.com/KubeOperator/ekko/internal/api/v1/role"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	"github.com/KubeOperator/ekko/internal/api/v1/user"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	v1RoleService "github.com/KubeOperator/ekko/internal/service/v1/role"
	v1RoleBindingService "github.com/KubeOperator/ekko/internal/service/v1/rolebinding"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
	"strings"
)

func authHandler() iris.Handler {
	return func(ctx *context.Context) {
		p := sessions.Get(ctx).Get("profile")
		if p == nil {
			ctx.Values().Set("message", "please login")
			ctx.StopWithStatus(iris.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}

func pageHandler() iris.Handler {
	return func(ctx *context.Context) {
		if ctx.URLParamExists(pkgV1.PageSize) && ctx.URLParamExists(pkgV1.PageNum) {
			pageNum, err := ctx.URLParamInt(pkgV1.PageNum)
			if err != nil {
				ctx.Values().Set("message", fmt.Sprintf("page num format err %s", err.Error()))
				ctx.StopWithStatus(iris.StatusBadRequest)
				return
			}
			pageSize, err := ctx.URLParamInt(pkgV1.PageSize)
			if err != nil {
				ctx.Values().Set("message", fmt.Sprintf("page size format err %s", err.Error()))
				ctx.StopWithStatus(iris.StatusBadRequest)
				return
			}
			ctx.Values().Set(pkgV1.PageNum, pageNum)
			ctx.Values().Set(pkgV1.PageSize, pageSize)
		}
		ctx.Next()
	}
}

func roleHandler() iris.Handler {
	return func(ctx *context.Context) {
		// 查询当前用户的角色
		// 查询角色的 rolebinding 获取 roles
		p := sessions.Get(ctx).Get("profile")
		u := p.(session.UserProfile)
		roleBindingService := v1RoleBindingService.NewService()
		rbs, err := roleBindingService.GetRoleBindingBySubject(v1Role.Subject{
			Kind: "User",
			Name: u.Name,
		})
		if err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				ctx.StatusCode(iris.StatusInternalServerError)
				ctx.Values().Set("message", err.Error())
				return
			}
		}
		roleNameHash := map[string]struct{}{}
		for i := range rbs {
			roleName := rbs[i].RoleRef
			roleNameHash[roleName] = struct{}{}
		}
		var roleNames []string
		for key := range roleNameHash {
			roleNames = append(roleNames, key)
		}
		ctx.Values().Set("roles", roleNames)
		ctx.Next()
	}
}

func roleAccessHandler() iris.Handler {
	return func(ctx *context.Context) {
		// 查询角色的 resources
		// 通过api resource 过滤出来资源主体,method 过滤操作
		p := sessions.Get(ctx).Get("profile")
		u := p.(session.UserProfile)
		roleNames := ctx.Values().Get("roles")
		roleService := v1RoleService.NewService()
		rs, _, err := roleService.Search(0, 0, pkgV1.Conditions{
			"Name": pkgV1.Condition{
				Field:    "Name",
				Operator: "in",
				Value:    roleNames,
			},
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		requestResource := func(path string) string {
			ss := strings.Split(path, "/")
			return strings.ToLower(ss[2])
		}(ctx.Request().URL.Path)
		requestMethod := strings.ToLower(ctx.Request().Method)

		resourceMatch := false
		methodMatch := false
		for i := range rs {
			for j := range rs[i].Rules {
				for k := range rs[i].Rules[j].Resource {
					if rs[i].Rules[j].Resource[k] == requestResource || rs[i].Rules[j].Resource[k] == "*" {
						resourceMatch = true
						for x := range rs[i].Rules[j].Verbs {
							if rs[i].Rules[j].Verbs[x] == requestMethod || rs[i].Rules[j].Verbs[x] == "*" {
								methodMatch = true
								break
							}
						}
					}
				}
			}
		}
		if !resourceMatch {
			ctx.StopWithStatus(iris.StatusForbidden)
			ctx.Values().Set("message", fmt.Sprintf("user %s has can't access  resource %s ", u.Name, requestResource))
			return
		}
		if !methodMatch {
			ctx.StopWithStatus(iris.StatusForbidden)
			ctx.Values().Set("message", fmt.Sprintf("user %s wants to use this verb %s", u.Name, strings.ToLower(ctx.Request().Method)))
			return
		}
		ctx.Next()
	}
}

func AddV1Route(app iris.Party) {
	v1Party := app.Party("/v1")
	v1Party.Use(pageHandler())
	session.Install(v1Party)
	authParty := v1Party.Party("")
	authParty.Use(authHandler())
	authParty.Use(roleHandler())
	authParty.Use(roleAccessHandler())
	user.Install(authParty)
	cluster.Install(v1Party)
	group.Install(authParty)
	role.Install(authParty)
}
