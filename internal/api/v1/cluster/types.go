package cluster

type Import struct {
	Name      string `json:"name"`
	ApiServer string `json:"apiServer"`
	Router    string `json:"router"`
	Token     string `json:"token"`
}
