package group

import (
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
)

type Group struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
}

type Binding struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	GroupRef     string `json:"groupRef" storm:"index"`
	UserRef      string `json:"userRef"`
}
