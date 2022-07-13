package config

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type Config struct {
	v1.BaseModel
	v1.Metadata
	Spec Spec `json:"spec"`
}
type Spec struct {
	Server  ServerConfig  `json:"server"`
	DB      DBConfig      `json:"db"`
	Session SessionConfig `json:"session"`
	Logger  LoggerConfig  `json:"logger"`
	AppId   string        `json:"appId"`
}

type ServerConfig struct {
	Bind BindConfig `json:"bind"`
	SSL  SSLConfig  `json:"ssl"`
}

type BindConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type SSLConfig struct {
	Enable         bool   `json:"enable"`
	Certificate    string `json:"certificate"`
	CertificateKey string `json:"certificateKey"`
}

type LoggerConfig struct {
	Level string
}

type DBConfig struct {
	Path string `json:"path"`
}

type SessionConfig struct {
	Expires int `json:"expires"`
}
