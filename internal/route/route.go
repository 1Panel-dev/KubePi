package route

import (
	v1 "github.com/1Panel-dev/KubePi/internal/api/v1"
	"github.com/kataras/iris/v12"
)

func InitRoute(party iris.Party) {
	apiParty := party.Party("/api")
	v1.AddV1Route(apiParty)
	//ws.AddWebSocketRoute(apiParty)
	//terminal.AddWebSocketRoute(apiParty)
}
