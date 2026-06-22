package cluster

import (
	"github.com/1Panel-dev/KubePi/internal/api/v1/session"
	"github.com/1Panel-dev/KubePi/internal/service/v1/clusteraccess"
	"github.com/1Panel-dev/KubePi/pkg/logging"
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
		previous := false
		/*是否显示日志时间*/
		timestamps := false
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
		if ctx.URLParamExists("timestamps") {
			p, err := ctx.URLParamBool("timestamps")
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				return
			}
			timestamps = p
		}

		sessionId, err := logging.GenLoggingSessionId()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
		}
		u := ctx.Values().Get("profile")
		profile := u.(session.UserProfile)
		client, _, _, err := clusteraccess.ClientForUser(clusterName, clusteraccess.User{
			Name:            profile.Name,
			IsAdministrator: profile.IsAdministrator,
		})
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.Values().Set("message", err)
			return
		}
		logging.LogSessions.Set(sessionId, logging.LogSession{
			Id:    sessionId,
			Bound: make(chan error),
		})
		go logging.WaitForLoggingStream(client, namespace, podName, containerName, tailLines, follow, previous, timestamps, sessionId)
		ctx.Values().Set("data", TerminalResponse{ID: sessionId})
	}
}
