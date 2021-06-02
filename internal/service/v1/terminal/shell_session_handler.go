package shell

import (
	"github.com/kataras/iris/v12/context"
)

type TerminalResponse struct {
	ID string `json:"id"`
}

func ExecShellHandler(ctx *context.Context) {
	//namespace := ctx.URLParam("namespace")
	//podName := ctx.URLParam("podName")
	//containerName := ctx.URLParam("containerName")
	//
	//sessionID, err := GenTerminalSessionId()
	//if err != nil {
	//	ctx.StatusCode(http.StatusInternalServerError)
	//	_, _ = ctx.JSON(map[string]interface{}{
	//		"msg": err.Error(),
	//	})
	//}
	//ekkoConfig := config.GetConfig()
	//var cfg *rest.Config
	//if ekkoConfig.KubeConfig != "" {
	//	cs, err := clientcmd.BuildConfigFromFlags("", ekkoConfig.KubeConfig)
	//	if err != nil {
	//		ctx.StatusCode(http.StatusInternalServerError)
	//		_, _ = ctx.JSON(err.Error())
	//		return
	//	}
	//	cfg = cs
	//}
	//client, err := kubernetes.NewForConfig(cfg)
	//terminalSessions.Set(sessionID, TerminalSession{
	//	id:       sessionID,
	//	bound:    make(chan error),
	//	sizeChan: make(chan remotecommand.TerminalSize),
	//})
	//go WaitForTerminal(client, cfg, namespace, podName, containerName, sessionID)
	//ctx.StatusCode(http.StatusOK)
	//_, _ = ctx.JSON(TerminalResponse{ID: sessionID})
}
