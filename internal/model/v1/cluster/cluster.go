package cluster

import v1 "github.com/KubeOperator/ekko/internal/model/v1"

type Cluster struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	ApiServer    string `json:"apiServer"`
	Router       string `json:"router"`
	Token        string `json:"token"`
	Status       string `json:"status"`
}
