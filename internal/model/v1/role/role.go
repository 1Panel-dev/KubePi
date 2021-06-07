package role

import v1 "github.com/KubeOperator/ekko/internal/model/v1"

type PolicyRule struct {
	Resource      []string `json:"resource"`
	ResourceNames []string `json:"resourceNames"`
	Verbs         []string `json:"verbs"`
}

type Role struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
}

type Subject struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type Banding struct {
	v1.BaseModel
	v1.Metadata
	Subjects []Subject `json:"subjects"`
	RoleRef  string    `json:"roleRef" storm:"index"`
}
