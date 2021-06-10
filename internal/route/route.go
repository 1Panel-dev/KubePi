package route

import (
	v1 "github.com/KubeOperator/ekko/internal/api/v1"
	"github.com/KubeOperator/ekko/internal/route/proxy"
	"github.com/kataras/iris/v12"
)

func InitRoute(party iris.Party) {
	apiParty := party.Party("/api")
	v1.AddV1Route(apiParty)
	proxyParty := party.Party("/proxy")
	proxy.AddProxyRoute(proxyParty)
	//ws.AddWebSocketRoute(apiParty)
	//terminal.AddWebSocketRoute(apiParty)
}
