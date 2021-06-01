package user

import v1 "github.com/KubeOperator/ekko/pkg/model/v1"

type User struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Spec         Spec `json:"spec" storm:"inline"`
}

type Info struct {
	NickName string `json:"nickName" storm:"index"`
	Email    string `json:"email" storm:"index,unique"`
}

type Authenticate struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Spec struct {
	Info         Info         `json:"info" storm:"inline"`
	Authenticate Authenticate `json:"authenticate" storm:"inline"`
}
