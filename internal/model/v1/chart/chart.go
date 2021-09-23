package chart

import (
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
)

type Chart struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Type         string `json:"type"`
	Spec         Spec   `json:"spec"`
	Status       Status `json:"status"`
}

type Spec struct {
	GitBranch      string         `json:"gitBranch"`
	GitRepo        string         `json:"gitRepo"`
	Url            string         `json:"url"`
	Authentication Authentication `json:"authentication"`
}

type Authentication struct {
	Type       string `json:"type"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PublicKey  string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
}

type Status struct {
	Version string `json:"version"`
	Phase   string `json:"phase"`
	Message string `json:"message"`
}
