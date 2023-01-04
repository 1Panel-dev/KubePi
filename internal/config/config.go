package config

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/model/v1/config"
	"github.com/KubeOperator/kubepi/pkg/file"
	"github.com/coreos/etcd/pkg/fileutil"
	"github.com/spf13/viper"
	"math/big"
	"strconv"
)

const configNotFoundSkipErr = "config file not found in %s, skip"
const configReadErr = "can not read config file %s ,%s"
const configMergeErr = "can not merge config file, %s"

var configFilePaths = []string{
	"/etc/kubepi",
}

func ReadConfig(c *config.Config, path ...string) error {
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
		return err
	}
	str, err := json.Marshal(&configMap)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(str, &c); err != nil {
		return nil
	}
	if c.Spec.Jwt.Key == "" {
		v.Set("spec.jwt.key", generate(32))
		if err := v.WriteConfig(); err != nil {
			return err
		}
	}
	return nil
}

func generate(length int) string {
	const base = 36
	size := big.NewInt(base)
	n := make([]byte, length)
	for i := range n {
		c, _ := rand.Int(rand.Reader, size)
		n[i] = strconv.FormatInt(c.Int64(), base)[0]
	}
	return string(n)
}
