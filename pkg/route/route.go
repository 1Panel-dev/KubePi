package route

import (
	"embed"
	"github.com/KubeOperator/ekko/pkg/route/proxy"
	"github.com/KubeOperator/ekko/pkg/route/terminal"
	"github.com/KubeOperator/ekko/pkg/route/ws"
	"github.com/kataras/iris/v12"
	"net/http"
)

var EmbedWebDashboard embed.FS
var EmbedWebTerminal embed.FS

func InitRoute(app *iris.Application) {
	dashboardFS := iris.PrefixDir("web/dashboard", http.FS(EmbedWebDashboard))
	app.RegisterView(iris.HTML(dashboardFS, ".html"))
	app.HandleDir("/ui", dashboardFS)

	terminalFS := iris.PrefixDir("web/terminal", http.FS(EmbedWebTerminal))
	app.RegisterView(iris.HTML(terminalFS, ".html"))
	app.HandleDir("/terminal", terminalFS)

	apiParty := app.Party("/api")
	proxy.AddProxyRoute(apiParty)
	ws.AddWebSocketRoute(apiParty)
	terminal.AddWebSocketRoute(apiParty)
}
