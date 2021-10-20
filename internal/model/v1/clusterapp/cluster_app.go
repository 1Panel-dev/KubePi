package clusterapp

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type ClusterApp struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	AppName      string `json:"appName"`
	Repo         string `json:"repo"`
	Cluster      string `json:"cluster"`
}
