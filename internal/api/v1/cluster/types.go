package cluster

import v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"

type Cluster struct {
	v1Cluster.Cluster
	KeyDataStr  string `json:"keyDataStr"`
	CertDataStr string `json:"certDataStr"`
	CaDataStr   string `json:"caDataStr"`
}
