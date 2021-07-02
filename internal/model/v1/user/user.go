package user

import v1 "github.com/KubeOperator/ekko/internal/model/v1"

type User struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	NickName     string `json:"nickName" storm:"index"`
	Email        string `json:"email" storm:"unique"`
	Language     string `json:"language"`
	Spec         Spec   `json:"spec" storm:"inline"`
}

type Authenticate struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Spec struct {
	Authenticate Authenticate `json:"authenticate" storm:"inline"`
}
