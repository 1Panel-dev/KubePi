package cluster

import (
	kContext "context"
	"fmt"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/KubeOperator/ekko/internal/util/kubernetes"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Handler struct {
	clusterService cluster.Cluster
}

func NewHandler() *Handler {
	return &Handler{
		clusterService: cluster.NewClusterService(),
	}
}

func (h *Handler) Create() iris.Handler {
	return func(ctx *context.Context) {
		var clusterImport Import
		if err := ctx.ReadJSON(&clusterImport); err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", err.Error())
			return
		}
		c, err := h.clusterService.Get(clusterImport.Name)
		if err != nil && err != storm.ErrNotFound {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("error: %s", err.Error()))
			return
		}
		if c != nil && c.Name != "" {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("cluster name %s is alerady exist", clusterImport.Name))
			return
		}
		client, err := kubernetes.NewKubernetesClient(&kubernetes.Config{
			Host:  clusterImport.ApiServer,
			Token: clusterImport.Token,
		})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("import cluster filed: %s", err.Error()))
			return
		}
		_, err = client.CoreV1().Namespaces().List(kContext.TODO(), metav1.ListOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("import cluster failed: %s", err.Error()))
			return
		}
		cs := v1Cluster.Cluster{
			Metadata: v1.Metadata{
				Name: clusterImport.Name,
				UUID: uuid.New().String(),
			},
			Token:     clusterImport.Token,
			Router:    clusterImport.Router,
			ApiServer: clusterImport.ApiServer,
			Status:    "Running",
		}
		err = h.clusterService.Create(&cs)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.Values().Set("message", fmt.Sprintf("create cluster failed: %s", err.Error()))
			return
		}
		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("data", cs)
	}
}

func (h *Handler) ListAll() iris.Handler {
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

func (h *Handler) Delete() iris.Handler {
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
	sp.Post("", handler.Create())
	sp.Get("", handler.ListAll())
	sp.Delete("/{name:string}", handler.Delete())
}
