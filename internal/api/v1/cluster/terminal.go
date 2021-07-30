package cluster

import (
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/KubeOperator/ekko/pkg/kubernetes"
	"github.com/KubeOperator/ekko/pkg/terminal"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"k8s.io/client-go/tools/remotecommand"
	"net/http"
)

type TerminalResponse struct {
	ID string `json:"id"`
}

func (h *Handler) TerminalSessionHandler() iris.Handler {
	return func(ctx *context.Context) {
		namespace := ctx.URLParam("namespace")
		podName := ctx.URLParam("podName")
		containerName := ctx.URLParam("containerName")

		sessionID, err := terminal.GenTerminalSessionId()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		clusterName := ctx.Params().GetString("name")
		c, err := h.clusterService.Get(clusterName, common.DBOptions{})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		k := kubernetes.NewKubernetes(*c)
		conf := k.Config()
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		terminal.TerminalSessions.Set(sessionID, terminal.TerminalSession{
			Id:       sessionID,
			Bound:    make(chan error),
			SizeChan: make(chan remotecommand.TerminalSize),
		})
		go terminal.WaitForTerminal(client, conf, namespace, podName, containerName, sessionID)
		ctx.StatusCode(http.StatusOK)
		resp := TerminalResponse{ID: sessionID}
		ctx.Values().Set("data", resp)
	}
}
