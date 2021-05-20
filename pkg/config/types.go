package config

type EkkoConfig struct {
	KubeApiServerConfig
	BindConfig
}

type KubeApiServerConfig struct {
	ApiServer  string
	Token      string
	KubeConfig string
}

type BindConfig struct {
	Host string
	Port int
}
