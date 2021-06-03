package v1

import (
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/cluster"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	"github.com/KubeOperator/ekko/internal/api/v1/user"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/sessions"
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

func AddV1Route(app iris.Party) {
	v1Party := app.Party("/v1")
	v1Party.Use(pageHandler())
	session.Install(v1Party)
	authParty := v1Party.Party("")
	authParty.Use(authHandler())
	user.Install(authParty)
	cluster.Install(v1Party)
}
