package repos

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	v2Url         = "api/v2.0"
	v1Url         = "api"
	projectUrl    = "/projects"
	repositoryUrl = "repositories"
	artifactUrl   = "/artifacts"
	tagUrl        = "tags"
)

func NewHarborClient(endpoint, username, password, version string) *harborClient {
	return &harborClient{
		EndPoint: endpoint,
		HttpClient: &HttpClient{
			Username: username,
			Password: password,
			Host:     endpoint,
		},
		Version: version,
	}
}

type harborClient struct {
	EndPoint   string
	HttpClient *HttpClient
	Version    string
}

type artifacts struct {
	Type string
	Tags []NameResult
}

func (c *harborClient) ListRepos() (names []string, err error) {

	result, err1 := c.HttpClient.GetNameResult(getProjectUrl(c.Version))
	if err1 != nil {
		err = err1
		return
	}
	for _, r := range result {
		names = append(names, r.Name)
	}
	return
}

func (c *harborClient) ListImages(project string) (images []string, err error) {

	if c.Version == "v2" {
		repoUrl := fmt.Sprintf("%s/%s/%s", getProjectUrl(c.Version), project, repositoryUrl)
		result, err1 := c.HttpClient.GetNameResult(repoUrl)
		if err1 != nil {
			err = err1
			return
		}
		for _, r := range result {
			repoName := strings.Replace(r.Name, project+"/", "", -1)
			body, _, err2 := c.HttpClient.Get(fmt.Sprintf("%s/%s/%s", repoUrl, repoName, artifactUrl))
			if err2 != nil {
				err = err2
				return
			}
			var artifacts []artifacts
			if err = json.Unmarshal(body, &artifacts); err != nil {
				return
			}
			for _, art := range artifacts {
				for _, tag := range art.Tags {
					images = append(images, r.Name+":"+tag.Name)
				}
			}
		}
	} else {

		result, err1 := c.HttpClient.GetNameResult(getProjectUrl(c.Version))
		if err1 != nil {
			err = err1
			return
		}
		var repoUrl string
		for _, value := range result {
			if value.Name == project {
				repoUrl = fmt.Sprintf("%s/%s%s", v1Url, repositoryUrl, "?project_id="+strconv.Itoa(value.ProjectID))
				break
			}
		}
		repos, err1 := c.HttpClient.GetNameResult(repoUrl)
		if err1 != nil {
			err = err1
			return
		}
		for _, repo := range repos {
			tagUrl := fmt.Sprintf("%s/%s/%s/%s", v1Url, repositoryUrl, repo.Name, tagUrl)
			tags, err1 := c.HttpClient.GetNameResult(tagUrl)
			if err1 != nil {
				err = err1
				return
			}
			for _, tag := range tags {
				images = append(images, repo.Name+":"+tag.Name)
			}
		}
	}
	return
}

func getProjectUrl(version string) string {
	if version == "v1" {
		return v1Url + projectUrl
	} else {
		return v2Url + projectUrl
	}
}
