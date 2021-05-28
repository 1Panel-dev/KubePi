package v1

import "time"

type BaseModel struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	CreateAt   time.Time
	UpdateAt   time.Time
}

type Metadata struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
