package repos

import (
	"encoding/json"
)

const (
	BaseUrl = "service/rest"
	RepoUrl = BaseUrl + "/v1/repositories"
)

func NewNexusClient(endpoint, username, password, downloadUrl string) *nexusClient {
	return &nexusClient{
		Username:    username,
		EndPoint:    endpoint,
		Password:    password,
		DownloadUrl: downloadUrl,
		HttpClient: &HttpClient{
			Username: username,
			Password: password,
			Host:     endpoint,
		},
	}
}

type nexusClient struct {
	Username    string
	Password    string
	EndPoint    string
	DownloadUrl string
	HttpClient  *HttpClient
}

type NameResult struct {
	Name string
}

func (c *nexusClient) Auth() error {
	return nil
}

func (c *nexusClient) ListRepos() (names []string, err error) {
	body, _, err1 := c.HttpClient.Get(RepoUrl)
	if err1 != nil {
		err = err1
		return
	}
	var result []NameResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	for _, r := range result {
		names = append(names, r.Name)
	}
	return
}
