package main

import (
	"github.com/KubeOperator/ekko/cmd/params"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/KubeOperator/ekko/pkg/server"
	"github.com/spf13/cobra"
)

var (
	kubeApiServerUrl string
	accessToken      string
	kubeConfigPath   string

	bindHost string
	bindPort int
)

func init() {
	rootCmd.Flags().StringVar(&kubeApiServerUrl, "kube-api-server-url", "http://127.0.0.1:8080", "kubernetes api server url")
	rootCmd.Flags().StringVar(&accessToken, "kube-access-token", "", "service account token")
	rootCmd.Flags().StringVar(&kubeConfigPath, "kube-config-path", "~/.kube/config", "kube config file path")
	rootCmd.Flags().StringVar(&bindHost, "bind-host", "0.0.0.0", "ekko bind address")
	rootCmd.Flags().IntVar(&bindPort, "bind-port", 2019, "ekko bind port")
}

var rootCmd = &cobra.Command{
	Use:   "ekko",
	Short: "A dashboard for kubernetes",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 验证参数合法性
		validator := params.NewParamValidator([]params.ParamValidateItem{
			{
				Value:        kubeApiServerUrl,
				ValidateFunc: params.KubeApiServerUrlValidateFunc,
			},
		})
		if err := validator.Check(); err != nil {
			return err
		}
		// 构建config对象
		c := config.EkkoConfig{
			KubeApiServerConfig: config.KubeApiServerConfig{
				ApiServer:  kubeApiServerUrl,
				Token:      accessToken,
				KubeConfig: kubeConfigPath,
			},
			BindConfig: config.BindConfig{
				Host: bindHost,
				Port: bindPort,
			},
		}
		config.SetConfig(c)
		// 启动 server
		return server.StartServer()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
