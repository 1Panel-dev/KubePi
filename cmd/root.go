package cmd

import (
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
		s := server.NewEkkoSerer()
		return s.Listen()
	},
}
