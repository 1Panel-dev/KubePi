package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"time"
)

const (
	clusterStatusInitializing = "Initializing"
	clusterStatusFailed       = "Failed"
	clusterStatusCompleted    = "Completed"
	clusterStatusSaved        = "Saved"
)

type Cluster struct {
	v1Cluster.Cluster
	KeyDataStr  string `json:"keyDataStr"`
	CertDataStr string `json:"certDataStr"`
	CaDataStr   string `json:"caDataStr"`
}

type NamespaceRoles struct {
	Namespace string   `json:"namespace"`
	Roles     []string `json:"roles"`
}

type Member struct {
	Name           string           `json:"name"`
	ClusterRoles   []string         `json:"clusterRoles"`
	BindingName    string           `json:"bindingName"`
	CreateAt       time.Time        `json:"createAt"`
	NamespaceRoles []NamespaceRoles `json:"namespaceRoles"`
}
