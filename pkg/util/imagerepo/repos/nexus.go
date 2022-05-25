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
	Items             []NameResult
	ContinuationToken string `json:"continuationToken"`
}

func (c *nexusClient) Auth() error {
	return nil
}

func (c *nexusClient) ListRepos(request ProjectRequest) (names []string, err error) {
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

func (c *nexusClient) ListImages(request RepoRequest) (response RepoResponse, err error) {
	var url string
	url = fmt.Sprintf("%s%s", ComponentUrl, request.Repo)
	if request.ContinueToken != "" {
		url += "&&continuationToken=" + request.ContinueToken
	}

	body, _, err1 := c.HttpClient.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	var result ItemResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	for _, r := range result.Items {
		response.Items = append(response.Items, r.Name+":"+r.Version)
	}
	response.ContinueToken = result.ContinuationToken
	return
}

func (c *nexusClient) ListImagesWithoutPage(repository string) (images []string, err error) {
	continuationToken := ""
	for {
		token := ""
		if continuationToken != "" {
			token = fmt.Sprintf("&&continuationToken=%s", continuationToken)
		}
		body, _, err1 := c.HttpClient.Get(fmt.Sprintf("%s%s", ComponentUrl, repository) + token)
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
		if result.ContinuationToken == "" {
			break
		}
		continuationToken = result.ContinuationToken
	}
	return
}
