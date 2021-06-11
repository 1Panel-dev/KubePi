package role

import v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"

type Role struct {
	v1Role.Role
	TemplateRef string `json:"templateRef"`
}
