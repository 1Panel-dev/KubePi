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

type harborBody struct {
	Name          string
	Version       string
	Format        string
	ProjectID     int `json:"project_id"`
	ArtifactCount int `json:"artifact_count"`
}

func (c *harborClient) ListRepos(request ProjectRequest) (names []string, err error) {

	projectUrl := fmt.Sprintf("%s?page=%d&page_size=%d", getProjectUrl(c.Version), request.Page, request.Limit)
	if request.Name != "" {
		projectUrl = projectUrl + "&q=name=~" + request.Name
	}
	body, _, err := c.HttpClient.Get(projectUrl)
	if err != nil {
		return
	}
	var projects []harborBody
	if err = json.Unmarshal(body, &projects); err != nil {
		return
	}
	for _, r := range projects {
		names = append(names, r.Name)
	}
	return
}

func (c *harborClient) ListImages(request RepoRequest) (response RepoResponse, err error) {
	project := request.Repo
	startCount := (request.Page - 1) * request.Limit

	if c.Version == "v2" {
		start := true
		p := 1
		limit := 100
		artifactCount := 0
		repoArCount := 0
		artifactStart := 0
		var repos []harborBody
		for start {
			repoUrl := fmt.Sprintf("%s/%s/%s?page=%d&&page_size=%d", getProjectUrl(c.Version), project, repositoryUrl, p, limit)
			body, _, err1 := c.HttpClient.Get(repoUrl)
			if err1 != nil {
				err = err1
				return
			}
			var items []harborBody
			if err1 = json.Unmarshal(body, &items); err1 != nil {
				err = err1
				return
			}
			if len(items) < limit {
				start = false
			}
			for _, v := range items {
				if v.ArtifactCount == 0 {
					continue
				}
				artifactCount = artifactCount + v.ArtifactCount
				if artifactCount >= startCount {
					repoArCount = repoArCount + v.ArtifactCount
					if repoArCount >= request.Limit {
						start = false
						break
					}
					repos = append(repos, v)
					continue
				}
				if artifactCount >= startCount {
					artifactStart = artifactCount - startCount
				}
			}
			p++
		}
		repoUrl := fmt.Sprintf("%s/%s/%s", getProjectUrl(c.Version), project, repositoryUrl)
		var items []string
		for _, r := range repos {
			repoName := strings.Replace(r.Name, project+"/", "", -1)
			body, res, err2 := c.HttpClient.Get(fmt.Sprintf("%s/%s/%s?page=%d&&page_size=%d", repoUrl, repoName, artifactUrl, 1, r.ArtifactCount))
			if res != nil && res.StatusCode == 404 {
				continue
			}
			if err2 != nil {
				if strings.Contains(err2.Error(), "404") {
					continue
				}
				err = err2
				return
			}
			var artifacts []artifacts
			if err = json.Unmarshal(body, &artifacts); err != nil {
				return
			}
			for _, art := range artifacts {
				for _, tag := range art.Tags {
					items = append(items, r.Name+":"+tag.Name)
				}
			}
		}
		artifactEnd := artifactStart + limit
		if artifactStart+limit > len(items) {
			artifactEnd = len(items) - 1
		}
		if artifactStart > artifactEnd {
			response.Items = []string{}
			return
		}
		response.Items = items[artifactStart : artifactEnd+1]
		if len(response.Items) == request.Limit {
			response.ContinueToken = "continue"
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
				response.Items = append(response.Items, repo.Name+":"+tag.Name)
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
