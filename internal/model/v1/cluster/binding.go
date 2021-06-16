package cluster

import v1 "github.com/KubeOperator/ekko/internal/model/v1"

type Subject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Binding struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Subjects     []Subject `json:"subjects"`
	ClusterRef   string    `json:"clusterRef" storm:"index"`
}
