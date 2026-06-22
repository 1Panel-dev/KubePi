package config

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/1Panel-dev/KubePi/internal/model/v1/config"
	"github.com/1Panel-dev/KubePi/pkg/file"
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

	paths := append([]string{}, configFilePaths...)
	paths = append(paths, path...)
	var loadedConfigFile string
	for i := range paths {
		realDir := file.ReplaceHomeDir(paths[i])

		if _, err := os.Stat(realDir); os.IsNotExist(err) {
			fmt.Println(fmt.Sprintf(configNotFoundSkipErr, realDir))
			continue
		} else if err != nil {
			fmt.Println(fmt.Sprintf(configReadErr, realDir, err.Error()))
			continue
		}

		configFile, ok, err := resolveConfigFile(realDir)
		if err != nil {
			fmt.Println(fmt.Sprintf(configReadErr, realDir, err.Error()))
			continue
		}
		if !ok {
			fmt.Println(fmt.Sprintf(configNotFoundSkipErr, realDir))
			continue
		}

		v.SetConfigFile(configFile)
		if err := v.MergeInConfig(); err != nil {
			fmt.Println(fmt.Sprintf(configMergeErr, paths))
			continue
		}
		loadedConfigFile = configFile

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
		c.Spec.Jwt.Key = generate(32)
		if loadedConfigFile != "" {
			v.Set("spec.jwt.key", c.Spec.Jwt.Key)
			v.SetConfigFile(loadedConfigFile)
			if err := v.WriteConfig(); err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveConfigFile(path string) (string, bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", false, err
	}
	if !info.IsDir() {
		return path, true, nil
	}

	for _, name := range []string{"app.yml", "app.yaml"} {
		configFile := filepath.Join(path, name)
		if _, err := os.Stat(configFile); err == nil {
			return configFile, true, nil
		} else if !os.IsNotExist(err) {
			return "", false, err
		}
	}
	return "", false, nil
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
