package main

import (
	"embed"
	"runtime"

	_ "github.com/KubeOperator/kubepi/cmd/server/docs"
	"github.com/KubeOperator/kubepi/internal/route"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/pkg/network/ip"
	"github.com/spf13/cobra"
)

//go:generate swag init

//swag init -g "cmd/server/main.go" -o "cmd/server/docs"

var (
	configPath     string
	serverBindHost string
	serverBindPort int
)

//go:embed web/kubepi
var embedWebKubePi embed.FS

//go:embed web/dashboard
var embedWebDashboard embed.FS

//go:embed web/terminal
var embedWebTerminal embed.FS

//go:embed script/darwin/init-kube.sh
var webkubectlEntrypointDarwin string

//go:embed script/linux/init-kube.sh
var webkubectlEntrypointLinux string

//go:embed helper/ip/qqwry.dat
var IpCommonDictionary []byte

func init() {
	RootCmd.Flags().StringVar(&serverBindHost, "server-bind-host", "", "kubepi bind address")
	RootCmd.Flags().IntVar(&serverBindPort, "server-bind-port", 0, "kubepi bind port")
	RootCmd.Flags().StringVarP(&configPath, "config-path", "c", "", "config file path")
}

var RootCmd = &cobra.Command{
	Use:   "kubepi-server",
	Short: "A dashboard for kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.EmbedWebDashboard = embedWebDashboard
		server.EmbedWebTerminal = embedWebTerminal
		server.EmbedWebKubePi = embedWebKubePi
		if runtime.GOOS == "darwin" {
			server.WebkubectlEntrypoint = webkubectlEntrypointDarwin
		} else {
			server.WebkubectlEntrypoint = webkubectlEntrypointLinux
		}
		ip.IpCommonDictionary = IpCommonDictionary
		return server.Listen(route.InitRoute,
			server.WithCustomConfigFilePath(configPath),
			server.WithServerBindHost(serverBindHost),
			server.WithServerBindPort(serverBindPort))
	},
}

// @title KubePi Restful API
// @version.go 1.0
// @termsOfService http://kubeoperator.io
// @contact.name Fit2cloud Support
// @contact.url https://www.fit2cloud.com
// @contact.email support@fit2cloud.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /kubepi/api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}
