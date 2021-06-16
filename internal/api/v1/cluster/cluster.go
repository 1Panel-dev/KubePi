package cluster

import (
	"fmt"
	"github.com/KubeOperator/ekko/internal/api/v1/session"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

type Handler struct {
	clusterService cluster.Cluster
}

func NewHandler() *Handler {
	return &Handler{
		clusterService: cluster.NewClusterService(),
	}
}

func (h *Handler) CreateCluster() iris.Handler {
	return func(ctx *context.Context) {
		var req Cluster
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		if req.CaDataStr != "" {
			req.CaCertificate.CertData = []byte(req.CaDataStr)
		}
		if req.Spec.Authentication.Mode == "certificate" {
			req.Spec.Authentication.Certificate.CertData = []byte(req.CertDataStr)
			req.Spec.Authentication.Certificate.KeyData = []byte(req.KeyDataStr)
		}

		client := kubernetes.NewKubernetes(req.Cluster)
		if err := client.Ping(); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		v, _ := client.Version()
		req.Status.Version = v.GitVersion
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		req.CreatedBy = profile.Name
		if err := h.clusterService.Create(&req.Cluster, common.DBOptions{}); err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
			return
		}
		ctx.Values().Set("data", &req)
	}
}

func (h *Handler) SearchClusters() iris.Handler {
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
		clusters, total, err := h.clusterService.Search(pageNum, pageSize, conditions, common.DBOptions{})
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err.Error())
		}
		ctx.Values().Set("data", pkgV1.Page{Items: clusters, Total: total})
	}
}

func (h *Handler) ListClusters() iris.Handler {
	return func(ctx *context.Context) {
		var clusters []v1Cluster.Cluster
		clusters, err := h.clusterService.All()
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("get clusters failed: %s", err.Error()))
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", clusters)
	}
}

func (h *Handler) DeleteCluster() iris.Handler {
	return func(ctx *context.Context) {
		name := ctx.Params().GetString("name")
		err := h.clusterService.Delete(name)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("delete cluster failed: %s", err.Error()))
			return
		}
		ctx.StatusCode(iris.StatusOK)
	}
}

func Install(parent iris.Party) {
	handler := NewHandler()
	sp := parent.Party("/clusters")
	sp.Post("", handler.CreateCluster())
	sp.Get("", handler.ListClusters())
	sp.Delete("/:name", handler.DeleteCluster())
	sp.Post("/search", handler.SearchClusters())
}
