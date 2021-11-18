package repos

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	baseUrl       = "api/v2.0"
	projectUrl    = baseUrl + "/projects"
	repositoryUrl = "/repositories"
	artifactUrl = "/artifacts"
)

func NewHarborClient(endpoint, username, password string) *harborClient {
	return &harborClient{
		EndPoint: endpoint,
		HttpClient: &HttpClient{
			Username: username,
			Password: password,
			Host:     endpoint,
		},
	}
}

type harborClient struct {
	EndPoint   string
	HttpClient *HttpClient
}

type  artifacts struct {
	Type string
	Tags []NameResult
}

func (c *harborClient) ListRepos() (names []string, err error) {
	body, _, err1 := c.HttpClient.Get(projectUrl)
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

func (c *harborClient) ListImages(project string) (images []string, err error) {
	repoUrl := fmt.Sprintf("%s/%s/%s", projectUrl, project, repositoryUrl)
	body, _, err1 := c.HttpClient.Get(repoUrl)
	if err1 != nil {
		err = err1
		return
	}
	var result []NameResult
	if err = json.Unmarshal(body, &result); err != nil {
		return
	}
	for _, r := range result {
		repoName := strings.Replace(r.Name,project+"/","",-1)
		body, _, err2 := c.HttpClient.Get(fmt.Sprintf("%s/%s/%s", repoUrl, repoName, artifactUrl))
		if err2 != nil {
			err = err2
			return
		}
		var artifacts []artifacts
		if err = json.Unmarshal(body, &artifacts); err != nil {
			return
		}
		for _,art := range artifacts {
			for _,tag := range art.Tags {
				images = append(images, r.Name+":"+tag.Name)
			}
		}
	}
	return
}
