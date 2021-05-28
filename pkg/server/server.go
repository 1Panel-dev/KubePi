package server

import (
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/KubeOperator/ekko/pkg/db"
	"github.com/KubeOperator/ekko/pkg/logger"
	"github.com/KubeOperator/ekko/pkg/route"
	"github.com/kataras/iris/v12"
)

type Phase interface {
	Name() string
	Init() error
}

var startPhase = []Phase{
	&db.InitDbPhase{},
}

func StartServer() error {
	config.Init()
	logger.Init()
	c := config.ReadConfig()
	l := logger.GetLogger()
	for i := range startPhase {
		l.Infof("start phase [%s]", startPhase[i].Name())
		if err := startPhase[i].Init(); err != nil {
			l.Fatalf("start phase [%s] error %s", startPhase[i].Name(), err.Error())
		}
	}
	app := iris.New()
	route.InitRoute(app)
	return app.Run(iris.Addr(fmt.Sprintf("%s:%d", c.Spec.Server.Bind.Host, c.Spec.Server.Bind.Port)))
}
