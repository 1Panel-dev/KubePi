package ldap

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type Ldap struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Address      string `json:"address"`
	Port         string `json:"port"`
	Dn           string `json:"dn"`
	Filter       string `json:"filter"`
}
