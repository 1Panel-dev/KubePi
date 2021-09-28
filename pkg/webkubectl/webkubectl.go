package webkubectl

import (
	"context"
	"github.com/KubeOperator/webkubectl/gotty/backend/localcommand"
	"github.com/KubeOperator/webkubectl/gotty/server"
	"github.com/KubeOperator/webkubectl/gotty/utils"
	"os"
)

type Interface interface {
	Run() error
}

func NewWebkubectl() Interface {
	return &webkubectl{}
}

type webkubectl struct {
}

func (w *webkubectl) Run() error {
	appOptions := &server.Options{}
	if err := utils.ApplyDefaultValues(appOptions); err != nil {
		return err
	}
	backendOptions := &localcommand.Options{}
	if err := utils.ApplyDefaultValues(backendOptions); err != nil {
		return err
	}
	redisOptions := &server.RedisOptions{}
	if err := utils.ApplyDefaultValues(redisOptions); err != nil {
		return err
	}
	args := []string{"bash"}
	hostname, _ := os.Hostname()
	appOptions.TitleVariables = map[string]interface{}{
		"command":  args[0],
		"argv":     args[1:],
		"hostname": hostname,
	}
	factory, err := localcommand.NewFactory(args[0], args[1:], backendOptions)
	if err != nil {
		return err
	}
	srv, err := server.New(factory, appOptions, redisOptions)
	if err != nil {
		return err
	}
	return srv.Run(context.Background())
}
