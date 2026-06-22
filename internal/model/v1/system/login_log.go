package system

import v1 "github.com/1Panel-dev/KubePi/internal/model/v1"

type LoginLog struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	UserName     string `json:"userName"`
	Ip           string `json:"ip"`
	City         string `json:"city"`
}
