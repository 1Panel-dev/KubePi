package session

import v1 "k8s.io/api/rbac/v1"

type LoginCredential struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordSetter struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type ProfileSetter struct {
	NickName string `json:"nickName"`
	Email    string `json:"email"`
}

type UserProfile struct {
	Name                string              `json:"name"`
	NickName            string              `json:"nickName"`
	Email               string              `json:"email"`
	ResourcePermissions map[string][]string `json:"resourcePermissions"`
}

type ClusterUserProfile struct {
	UserProfile
	ClusterRoles []v1.ClusterRole `json:"clusterRoles"`
}
