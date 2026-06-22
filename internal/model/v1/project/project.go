package project

import v1 "github.com/1Panel-dev/KubePi/internal/model/v1"

type Project struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Users        []string `json:"users"`
	ClusterRef   string   `json:"clusterRef" storm:"index"`
	ProjectName  string   `json:"projectName"`
	Certificate  []byte   `json:"certificate"`
}
