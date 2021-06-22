package group

import (
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"time"
)

type Member struct {
	Name     string    `json:"name"`
	Group    string    `json:"group"`
	CreateBy string    `json:"createBy"`
	CreateAt time.Time `json:"createAt"`
}

type Group struct {
	v1Group.Group
	Roles []string `json:"roles"`
}
