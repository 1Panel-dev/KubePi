package imagerepo

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type ImageRepo struct {
	v1.BaseModel   `storm:"inline"`
	v1.Metadata    `storm:"inline"`
	Type           string     `json:"type"`
	EndPoint       string     `json:"endPoint"`
	DownloadUrl    string     `json:"downloadUrl"`
	RepoName       string     `json:"repoName"`
	Credential     Credential `json:"credential"`
	Auth           bool       `json:"auth"`
	AllowAnonymous bool       `json:"allowAnonymous"`
	Version        string     `json:"version"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
