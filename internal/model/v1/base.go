package v1

import "time"

type BaseModel struct {
	ApiVersion string    `json:"apiVersion"`
	Kind       string    `json:"kind"`
	CreateAt   time.Time `json:"createAt" storm:"index"`
	UpdateAt   time.Time `json:"updateAt" storm:"index"`
	BuiltIn    bool      `json:"builtIn"`
	CreatedBy  string    `json:"createdBy"`
}

type Metadata struct {
	Name        string `json:"name" storm:"unique" `
	Description string `json:"description"`
	UUID        string `json:"uuid" storm:"id,index,unique"`
}
