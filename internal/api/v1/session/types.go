package session

import v1 "k8s.io/api/rbac/v1"

type LoginCredential struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	AuthMethod string `json:"authMethod"`
}
type MfaCredential struct {
	Username string `json:"username"`
	Secret   string `json:"secret"`
	Code     string `json:"code"`
}

type PasswordSetter struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type ProfileSetter struct {
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	Language string `json:"language"`
}

type UserProfile struct {
	Name                string              `json:"name"`
	NickName            string              `json:"nickName"`
	Email               string              `json:"email"`
	Language            string              `json:"language"`
	ResourcePermissions map[string][]string `json:"resourcePermissions"`
	IsAdministrator     bool                `json:"isAdministrator"`
	Mfa                 Mfa                 `json:"mfa"`
}

type ClusterUserProfile struct {
	UserProfile
	ClusterRoles []v1.ClusterRole `json:"clusterRoles"`
}

type Mfa struct {
	Enable   bool   `json:"enable"`
	Secret   string `json:"secret"`
	Approved bool   `json:"approved"`
}
