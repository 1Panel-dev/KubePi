package config

import (
	"encoding/json"
	"fmt"
	v1 "github.com/KubeOperator/ekko/pkg/model/v1"
	"github.com/KubeOperator/ekko/pkg/model/v1/config"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/viper"
)

var v *viper.Viper
var configData config.Config

const configErr = "can not parse config file , %s"
const configNotFoundSkipErr = "config file not found in %s, skip"
const configReadErr = "can not read config file %s ,%s"
const configMergeErr = "can not merge config file, %s"

var configFilePaths = []string{
	"/etc/ekko",
}

func Init(path ...string) {
	v = viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")

	for i := range path {
		configFilePaths = append(configFilePaths, path[i])
	}
	for i := range configFilePaths {
		if exists := fileutil.Exist(configFilePaths[i]); !exists {
			fmt.Println(fmt.Sprintf(configNotFoundSkipErr, configFilePaths[i]))
			continue
		}
		v.AddConfigPath(configFilePaths[i])
		if err := v.ReadInConfig(); err != nil {
			fmt.Println(fmt.Sprintf(configReadErr, configFilePaths[i], err.Error()))
			continue
		}
		if err := v.MergeInConfig(); err != nil {
			fmt.Println(fmt.Sprintf(configMergeErr, configFilePaths))
		}
	}

	var configMap map[string]interface{}
	if err := v.Unmarshal(&configMap); err != nil {
		configErrHandler(err)
	}
	str, err := json.Marshal(&configMap)
	if err != nil {
		configErrHandler(err)
	}
	configData = defaultConfig()
	if err := json.Unmarshal(str, &configData); err != nil {
		configErrHandler(err)
	}
}

func ReadConfig() config.Config {
	return configData
}

func GetConfigInstance() *config.Config {
	return &configData
}

func configErrHandler(err error) {
	panic(fmt.Sprintf(configErr, err.Error()))
}

func defaultConfig() config.Config {
	return config.Config{
		BaseModel: v1.BaseModel{
			ApiVersion: "v1",
			Kind:       "AppConfig",
		},
		Metadata: v1.Metadata{},
		Spec: config.Spec{
			Server: config.ServerConfig{
				Bind: config.BindConfig{
					Host: "0.0.0.0",
					Port: 2019,
				},
				SSL: config.SSLConfig{
					Enable: false,
				},
			},
			DB: config.DBConfig{
				Path: "/var/lib/ekko/db/ekko.db",
			},
		},
	}
}
