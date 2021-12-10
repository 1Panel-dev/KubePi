package repos

import (
	"encoding/json"
	"fmt"
)

const (
	BaseUrl      = "service/rest"
	RepoUrl      = BaseUrl + "/v1/repositories"
	ComponentUrl = BaseUrl + "/v1/components?repository="
)

func NewNexusClient(endpoint, username, password string) *nexusClient {
	return &nexusClient{
		Username: username,
		EndPoint: endpoint,
		Password: password,
		HttpClient: &HttpClient{
			Username: username,
			Password: password,
			Host:     endpoint,
		},
	}
}

type nexusClient struct {
	Username   string
	Password   string
	EndPoint   string
	HttpClient *HttpClient
}

type ItemResult struct {
	Items []NameResult
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
		if r.Format == "docker" {
			names = append(names, r.Name)
		}
	}
	return
}

func (c *nexusClient) ListImages(repository string) (images []string, err error) {
	body, _, err1 := c.HttpClient.Get(fmt.Sprintf("%s%s", ComponentUrl, repository))
	if err1 != nil {
		err = err1
		return
	}
	var result ItemResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	for _, r := range result.Items {
		images = append(images, r.Name+":"+r.Version)
	}
	return
}
