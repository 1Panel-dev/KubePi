package cluster

import (
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
)

type Cluster struct {
	v1.BaseModel  `storm:"inline"`
	v1.Metadata   `storm:"inline"`
	CaCertificate Certificate `json:"caCertificate" storm:"inline"`
	Spec          Spec        `json:"spec" storm:"inline"`
	PrivateKey    []byte      `json:"privateKey"`
	Status        Status      `json:"status" storm:"inline"`
	Labels        []string    `json:"labels"`
}

type Spec struct {
	Connect        Connect        `json:"connect" storm:"inline"`
	Authentication Authentication `json:"authentication" storm:"inline"`
	Local          bool           `json:"local"`
}

type Connect struct {
	Direction string  `json:"direction"`
	Forward   Forward `json:"forward" storm:"inline"`
}

type Forward struct {
	ApiServer string `json:"apiServer"`
	Proxy     Proxy  `json:"proxy"   storm:"inline"`
}

type Proxy struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Authentication struct {
	Mode              string      `json:"mode"`
	BearerToken       string      `json:"bearerToken"`
	Certificate       Certificate `json:"certificate" storm:"inline"`
	ConfigFileContent []byte      `json:"configFileContent"`
}

type Certificate struct {
	KeyData  []byte `json:"keyData"`
	CertData []byte `json:"certData"`
}

type Status struct {
	Version string `json:"version"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
}
