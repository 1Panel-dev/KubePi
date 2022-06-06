package docs

type UserCreate struct {
	Name         string       `json:"name"`
	NickName     string       `json:"nickName" storm:"index"`
	Email        string       `json:"email" storm:"unique"`
	IsAdmin      bool         `json:"isAdmin"`
	Authenticate Authenticate `json:"authenticate"`
	Mfa          Mfa          `json:"mfa"`
	Roles        []string     `json:"roles"`
}

type Authenticate struct {
	Password string `json:"password"`
}

type Mfa struct {
	Enable bool `json:"enable"`
}
