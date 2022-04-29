package repos

type Config struct {
	Type     string
	EndPoint string
	//DownloadUrl string
	Credential Credential
	Version    string
}

type Credential struct {
	Username string
	Password string
}

type RepoRequest struct {
	Repo          string `json:"repo"`
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	ContinueToken string `json:"continueToken"`
}

type ProjectRequest struct {
	Page          int    `json:"page"`
	Limit         int    `json:"limit"`
	Name          string `json:"name"`
	ContinueToken string `json:"continueToken"`
}

type RepoResponse struct {
	Items         []string `json:"items"`
	Total         int      `json:"total"`
	ContinueToken string   `json:"continueToken"`
}
