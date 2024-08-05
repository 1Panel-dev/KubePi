package cluster

import (
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/terminal"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"k8s.io/client-go/tools/remotecommand"
)

type TerminalResponse struct {
	ID string `json:"id"`
}

func (h *Handler) TerminalSessionHandler() iris.Handler {
	return func(ctx *context.Context) {
		namespace := ctx.URLParam("namespace")
		podName := ctx.URLParam("podName")
		containerName := ctx.URLParam("containerName")
		shell := ctx.URLParam("shell")

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
		k := kubernetes.NewKubernetes(c)
		conf, err := k.Config()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		client, err := k.Client()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		if shell == "" {
			shell = "sh"
		}
		terminal.TerminalSessions.Set(sessionID, terminal.TerminalSession{
			Id:       sessionID,
			Bound:    make(chan error),
			SizeChan: make(chan remotecommand.TerminalSize),
		})
		go terminal.WaitForTerminal(client, conf, namespace, podName, containerName, sessionID, shell)
		resp := TerminalResponse{ID: sessionID}
		ctx.Values().Set("data", resp)
	}
}

//node shell
func (h *Handler) NodeTerminalSessionHandler() iris.Handler {
	return func(ctx *context.Context) {
		nodeName := ctx.URLParam("nodeName")

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
		k := kubernetes.NewKubernetes(c)
		conf, err := k.Config()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
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
		go terminal.WaitForNodeShellTerminal(client, conf, nodeName, sessionID)
		resp := TerminalResponse{ID: sessionID}
		ctx.Values().Set("data", resp)
	}
}
