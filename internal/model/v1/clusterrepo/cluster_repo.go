package clusterrepo

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type ClusterRepo struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Cluster      string `json:"cluster"`
	Repo         string `json:"repo"`
}

