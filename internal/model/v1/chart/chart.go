package chart

type RepoCreate struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}
