package user

import v1 "github.com/KubeOperator/ekko/pkg/model/v1"

type User struct {
	v1.BaseModel
	v1.Metadata
}

type Info struct {
	NickName string `json:"nickName"`
	Email    string `json:"email"`
}

type Authenticate struct {
	Password string `json:"password"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
}

type Spec struct {
	Info         Info         `json:"info"`
	Authenticate Authenticate `json:"authenticate"`
}
