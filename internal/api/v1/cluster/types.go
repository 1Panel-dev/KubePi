package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"time"
)

type Cluster struct {
	v1Cluster.Cluster
	KeyDataStr  string `json:"keyDataStr"`
	CertDataStr string `json:"certDataStr"`
	CaDataStr   string `json:"caDataStr"`
}

type Member struct {
	Name         string    `json:"name"`
	Kind         string    `json:"kind"`
	ClusterRoles []string  `json:"clusterRoles"`
	CreateAt     time.Time `json:"createAt"`
	BindingName  string    `json:"bindingName"`
}
