package server

import (
	"fmt"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/KubeOperator/ekko/pkg/route"
	"github.com/kataras/iris/v12"
)

func StartServer() error {
	conf := config.GetConfig()
	app := iris.New()
	route.InitRoute(app)
	return app.Run(iris.Addr(fmt.Sprintf("%s:%d", conf.Host, conf.Port)))
}
