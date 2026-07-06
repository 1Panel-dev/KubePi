package cluster

import (
	"errors"

	"github.com/1Panel-dev/KubePi/internal/api/v1/session"
	"github.com/1Panel-dev/KubePi/internal/service/v1/clusteraccess"
	"github.com/1Panel-dev/KubePi/internal/service/v1/common"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func userAccessFromContext(ctx *context.Context) (clusteraccess.User, bool) {
	profile, ok := ctx.Values().Get("profile").(session.UserProfile)
	if !ok || profile.Name == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.Values().Set("message", "please login")
		return clusteraccess.User{}, false
	}
	return clusteraccess.User{Name: profile.Name, IsAdministrator: profile.IsAdministrator}, true
}

func writeClusterAccessError(ctx *context.Context, err error) {
	if errors.Is(err, clusteraccess.ErrClusterAccessDenied) {
		ctx.StatusCode(iris.StatusForbidden)
		ctx.Values().Set("message", "user can not access cluster")
		return
	}
	ctx.StatusCode(iris.StatusInternalServerError)
	ctx.Values().Set("message", err.Error())
}

func (h *Handler) requireClusterAccess() iris.Handler {
	return func(ctx *context.Context) {
		user, ok := userAccessFromContext(ctx)
		if !ok {
			return
		}
		clusterName := ctx.Params().GetString("name")
		if err := clusteraccess.CheckClusterAccess(clusterName, user, common.DBOptions{}); err != nil {
			writeClusterAccessError(ctx, err)
			return
		}
		ctx.Next()
	}
}

func (h *Handler) clusterAccessSet(profile session.UserProfile) (map[string]struct{}, error) {
	accessSet := map[string]struct{}{}
	if profile.IsAdministrator {
		return accessSet, nil
	}
	bindings, err := h.clusterBindingService.GetBindingsByUserName(profile.Name, common.DBOptions{})
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return accessSet, nil
		}
		return nil, err
	}
	for i := range bindings {
		accessSet[bindings[i].ClusterRef] = struct{}{}
	}
	return accessSet, nil
}

func canAccessCluster(profile session.UserProfile, accessSet map[string]struct{}, clusterName string) bool {
	if profile.IsAdministrator {
		return true
	}
	_, ok := accessSet[clusterName]
	return ok
}

func paginateClusterResults(result []Cluster, pageNum, pageSize int) []Cluster {
	if pageSize <= 0 {
		return result
	}
	if pageNum < 1 {
		pageNum = 1
	}
	start := (pageNum - 1) * pageSize
	if start >= len(result) {
		return []Cluster{}
	}
	end := start + pageSize
	if end > len(result) {
		end = len(result)
	}
	return result[start:end]
}
