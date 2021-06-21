package proxy

import (
	"github.com/KubeOperator/ekko/internal/service/v1/cluster"
	"github.com/kataras/iris/v12"
)

var (
	clusterService = cluster.NewService()
)

func AddProxyRoute(app iris.Party) {
	app.Any("/{cluster_name}/{p:path}", KubernetesClientProxy)
}
