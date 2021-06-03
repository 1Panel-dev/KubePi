package user

import (
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
)

type PolicyRule struct {
	Resource      []string
	ResourceNames []string
	Verbs         []string
}

type Role struct {
	v1.BaseModel
	v1.Metadata
}

type Subject struct {
	Kind string
	Name string
}

type RoleBanding struct {
	v1.BaseModel
	v1.Metadata
	Subjects []Subject
	RoleRef  string
}
