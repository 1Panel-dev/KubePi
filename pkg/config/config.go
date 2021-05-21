package config

var config EkkoConfig

func GetConfig() EkkoConfig {
	return config
}

func SetConfig(c EkkoConfig) {
	config = c
}
