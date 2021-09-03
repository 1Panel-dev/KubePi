package cluster

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type Binding struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	UserRef      string `json:"UserRef" storm:"inline"`
	ClusterRef   string `json:"clusterRef" storm:"index"`
	Certificate  []byte `json:"certificate"`
}
