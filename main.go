package main

import (
	"embed"
	"github.com/KubeOperator/ekko/cmd"
	"github.com/KubeOperator/ekko/pkg/server"
	"github.com/spf13/cobra"
)

//go:embed web/terminal
var embedWebTerminal embed.FS

//go:embed web/dashboard
var embedWebDashboard embed.FS

func main() {
	server.EmbedWebDashboard = embedWebDashboard
	server.EmbedWebTerminal = embedWebTerminal
	cobra.CheckErr(cmd.RootCmd.Execute())
}
