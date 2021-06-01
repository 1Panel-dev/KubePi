package v1

import (
	"github.com/KubeOperator/ekko/pkg/api/v1/session"
	"github.com/kataras/iris/v12"
)

func AddV1Route(app iris.Party) {
	v1Party := app.Party("/v1")
	session.Install(v1Party)
}
