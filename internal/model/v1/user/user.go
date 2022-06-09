package user

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type User struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	NickName     string       `json:"nickName" storm:"index"`
	Email        string       `json:"email" storm:"unique"`
	Language     string       `json:"language"`
	IsAdmin      bool         `json:"isAdmin"`
	Authenticate Authenticate `json:"authenticate"`
	Type         string       `json:"type"`
	Mfa          Mfa          `json:"mfa"`
}

type Authenticate struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type Mfa struct {
	Enable bool   `json:"enable"`
	Secret string `json:"secret"`
}

const (
	LDAP  = "LDAP"
	LOCAL = "LOCAL"
)

type ImportUser struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	NickName  string `json:"nickName"`
	Available bool   `json:"available"`
}

type ImportResult struct {
	Success  bool     `json:"success"`
	Failures []string `json:"failures"`
	Msg      string   `json:"msg"`
}
