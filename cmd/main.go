package main

import (
	"fmt"
	"github.com/KubeOperator/ekko/cmd/params"
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/KubeOperator/ekko/pkg/server"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"strings"
)

var (
	kubeApiServerUrl string
	accessToken      string
	kubeConfigPath   string

	bindHost string
	bindPort int
)

func init() {
	rootCmd.Flags().StringVar(&kubeApiServerUrl, "kube-api-server-url", "", "kubernetes api server url, eg:https://127.0.0.1:8443")
	rootCmd.Flags().StringVar(&accessToken, "kube-access-token", "", "service account token")
	rootCmd.Flags().StringVar(&kubeConfigPath, "kube-config-path", "${HOME_DIR}/.kube/config", "kube config file path")
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
		c := config.EkkoConfig{BindConfig: config.BindConfig{
			Host: bindHost,
			Port: bindPort,
		}}

		if (accessToken == "" && kubeApiServerUrl != "") || (accessToken != "" && kubeApiServerUrl == "") {
			return fmt.Errorf("kube-api-server-url and kube-access-token must be set together")
		}
		if accessToken != "" && kubeApiServerUrl != "" {
			fmt.Printf("use token to connect cluster: %s \n", kubeApiServerUrl)
			c.KubeApiServerConfig.ApiServer = kubeApiServerUrl
			c.KubeApiServerConfig.Token = accessToken
		}
		if accessToken == "" && kubeApiServerUrl == "" {
			if strings.Contains(kubeConfigPath, "${HOME_DIR}") || strings.Contains(kubeConfigPath, "~") {
				hd, _ := homedir.Dir()
				kubeConfigPath = strings.Replace(kubeConfigPath, "${HOME_DIR}", hd, -1)
				kubeConfigPath = strings.Replace(kubeConfigPath, "~", hd, -1)
			}
			fmt.Printf("use kube-config-file :%s to connect cluster \n", kubeConfigPath)
			if exists := fileutil.Exist(kubeConfigPath); !exists {
				return fmt.Errorf("%s,no valid connection credentials were found", kubeConfigPath)
			}
			c.KubeConfig = kubeConfigPath
		}
		config.SetConfig(c)
		// 启动 server
		return server.StartServer()
	},
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
