package terminal

import (
	shell "github.com/KubeOperator/ekko/internal/terminal"
	"github.com/kataras/iris/v12"
)

func AddWebSocketRoute(app iris.Party) {
	terminalParty := app.Party("/terminal")
	terminalParty.Get("/session", shell.ExecShellHandler)
}
