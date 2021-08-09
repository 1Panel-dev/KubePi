package config

import (
	"encoding/json"
	"fmt"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	"github.com/KubeOperator/ekko/internal/model/v1/config"
	"github.com/KubeOperator/ekko/pkg/file"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/viper"
)

const configNotFoundSkipErr = "config file not found in %s, skip"
const configReadErr = "can not read config file %s ,%s"
const configMergeErr = "can not merge config file, %s"

var configFilePaths = []string{
	"~/.ekko/conf",
}

func ReadConfig(path ...string) (config.Config, error) {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")

	for i := range path {
		configFilePaths = append(configFilePaths, path[i])
	}
	for i := range configFilePaths {
		realDir := file.ReplaceHomeDir(configFilePaths[i])

		if exists := fileutil.Exist(realDir); !exists {
			fmt.Println(fmt.Sprintf(configNotFoundSkipErr, realDir))
			continue
		}

		v.AddConfigPath(realDir)
		if err := v.ReadInConfig(); err != nil {
			fmt.Println(fmt.Sprintf(configReadErr, realDir, err.Error()))
			continue
		}
		if err := v.MergeInConfig(); err != nil {
			fmt.Println(fmt.Sprintf(configMergeErr, configFilePaths))
		}
	}

	var configMap map[string]interface{}
	if err := v.Unmarshal(&configMap); err != nil {
		return config.Config{}, err
	}
	str, err := json.Marshal(&configMap)
	if err != nil {
		return config.Config{}, err
	}
	c := defaultConfig()
	if err := json.Unmarshal(str, &c); err != nil {
		return config.Config{}, nil
	}
	return c, nil
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
				Path: "~/.ekko/db",
			},
			Logger: config.LoggerConfig{Level: "debug"},
		},
	}
}
