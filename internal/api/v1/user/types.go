package user

import v1User "github.com/1Panel-dev/KubePi/internal/model/v1/user"

type User struct {
	v1User.User
	Roles       []string `json:"roles"`
	OldPassword string   `json:"oldPassword"`
	Password    string   `json:"password"`
}
