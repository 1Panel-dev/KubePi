package session

type LoginCredential struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	Name     string `json:"name"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
}
