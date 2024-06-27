package cluster

import (
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/logging"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func (h *Handler) LoggingHandler() iris.Handler {
	return func(ctx *context.Context) {
		clusterName := ctx.Params().GetString("name")

		namespace := ctx.URLParam("namespace")
		podName := ctx.URLParam("podName")
		containerName := ctx.URLParam("containerName")
		tailLines := 100
		follow := false
/*是否查看上次失败的容器日志*/
		previous :=false
		if ctx.URLParamExists("tailLines") {
			lines, err := ctx.URLParamInt("tailLines")
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				return
			}
			tailLines = lines
		}
		if ctx.URLParamExists("follow") {
			f, err := ctx.URLParamBool("follow")
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				return
			}
			follow = f
		}
if ctx.URLParamExists("previous") {
			p, err := ctx.URLParamBool("previous")
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				return
			}
			previous = p
		}
		sessionId, err := logging.GenLoggingSessionId()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
		c, err := h.clusterService.Get(clusterName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		k := kubernetes.NewKubernetes(c)
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		logging.LogSessions.Set(sessionId, logging.LogSession{
			Id:    sessionId,
			Bound: make(chan error),
		})
		go logging.WaitForLoggingStream(client, namespace, podName, containerName, tailLines,  follow, previous, sessionId)
		ctx.Values().Set("data", TerminalResponse{ID: sessionId})
	}
}
