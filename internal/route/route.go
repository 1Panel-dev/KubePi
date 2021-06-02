package route

import (
	v1 "github.com/KubeOperator/ekko/internal/api/v1"
	"github.com/kataras/iris/v12"
)

func InitRoute(party iris.Party) {
	apiParty := party.Party("/api")
	v1.AddV1Route(apiParty)
	//proxy.AddProxyRoute(apiParty)
	//ws.AddWebSocketRoute(apiParty)
	//terminal.AddWebSocketRoute(apiParty)
}
