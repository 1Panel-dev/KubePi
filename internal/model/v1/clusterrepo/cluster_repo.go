package clusterrepo

import v1 "github.com/1Panel-dev/KubePi/internal/model/v1"

type ClusterRepo struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Cluster      string `json:"cluster"`
	Repo         string `json:"repo"`
}
