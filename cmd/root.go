package cmd

import (
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/KubeOperator/ekko/pkg/server"
	"github.com/spf13/cobra"
)

var (
	configPath     string
	serverBindHost string
	serverBindPort int
)

func init() {
	RootCmd.Flags().StringVar(&serverBindHost, "server-bind-host", "", "ekko bind address")
	RootCmd.Flags().IntVar(&serverBindPort, "server-bind-port", 0, "ekko bind port")
	RootCmd.Flags().StringVarP(&configPath, "config-path", "c", "", "config file path")
}

var RootCmd = &cobra.Command{
	Use:   "ekko",
	Short: "A dashboard for kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configPath != "" {
			config.Init(configPath)
		} else {
			config.Init()
		}
		c := config.GetConfigInstance()
		if serverBindHost != "" {
			c.Spec.Server.Bind.Host = serverBindHost
		}
		if serverBindPort != 0 {
			c.Spec.Server.Bind.Port = serverBindPort
		}
		// 启动 server
		return server.StartServer()
	},
}
