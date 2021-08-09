package main

import (
	"embed"

	"github.com/KubeOperator/ekko/internal/route"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/spf13/cobra"
)

var (
	configPath     string
	serverBindHost string
	serverBindPort int
)

//go:embed web/ekko
var embedWebEkko embed.FS

//go:embed web/dashboard
var embedWebDashboard embed.FS

//go:embed web/terminal
var embedWebTerminal embed.FS

func init() {
	RootCmd.Flags().StringVar(&serverBindHost, "server-bind-host", "", "ekko bind address")
	RootCmd.Flags().IntVar(&serverBindPort, "server-bind-port", 0, "ekko bind port")
	RootCmd.Flags().StringVarP(&configPath, "config-path", "c", "", "config file path")
}

var RootCmd = &cobra.Command{
	Use:   "ekko-server",
	Short: "A dashboard for kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.EmbedWebDashboard = embedWebDashboard
		server.EmbedWebTerminal = embedWebTerminal
		server.EmbedWebEkko = embedWebEkko
		return server.Listen(route.InitRoute)
	},
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
