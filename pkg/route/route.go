package route

import (
	"github.com/KubeOperator/ekko/pkg/route/proxy"
	"github.com/kataras/iris/v12"
)

func InitRoute(app *iris.Application) {
	app.Party("/proxy").Any("/{p:path}", proxy.KubernetesApiProxy)
}
