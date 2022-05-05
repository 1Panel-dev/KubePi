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
	TagsCount     int `json:"tags_count"`
}

func (c *harborClient) ListRepos(request ProjectRequest) (names []string, err error) {

	projectUrl := fmt.Sprintf("%s?page=%d&page_size=%d", getProjectUrl(c.Version), request.Page, request.Limit)
	if request.Name != "" {
		if c.Version == "v2" {
			projectUrl = projectUrl + "&q=name=~" + request.Name
		} else {
			projectUrl = projectUrl + "&name=" + request.Name
		}
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

const LIMIT = 100

func (c *harborClient) ListImages(request RepoRequest) (response RepoResponse, err error) {
	project := request.Repo
	startCount := (request.Page - 1) * request.Limit

	start := true
	p := 1
	artifactCount := 0
	repoArCount := 0
	artifactStart := 0
	var repos []harborBody
	for start {
		var repoUrl string
		if c.Version == "v2" {
			repoUrl = fmt.Sprintf("%s/%s/%s?page=%d&&page_size=%d", getProjectUrl(c.Version), project, repositoryUrl, p, LIMIT)
		} else {
			result, err1 := c.HttpClient.GetNameResult(getProjectUrl(c.Version))
			if err1 != nil {
				err = err1
				return
			}
			for _, value := range result {
				if value.Name == project {
					repoUrl = fmt.Sprintf("%s/%s%s", v1Url, repositoryUrl, "?project_id="+strconv.Itoa(value.ProjectID))
					break
				}
			}
		}

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
		if len(items) < LIMIT {
			start = false
		}
		for _, v := range items {

			count := v.ArtifactCount
			if c.Version == "v1" {
				count = v.TagsCount
			}

			if count == 0 {
				continue
			}
			artifactCount = artifactCount + count
			if artifactCount > startCount {
				repos = append(repos, v)
				repoArCount = repoArCount + count
				if repoArCount >= request.Limit {
					start = false
					break
				}
				continue
			}
			if artifactCount >= startCount {
				artifactStart = artifactCount - startCount
			}
		}
		p++
	}

	var repoUrl string
	if c.Version == "v2" {
		repoUrl = fmt.Sprintf("%s/%s/%s", getProjectUrl(c.Version), project, repositoryUrl)
	}
	var items []string
	for _, r := range repos {
		var tagsUrl string
		if c.Version == "v2" {
			repoName := strings.Replace(r.Name, project+"/", "", -1)
			tagsUrl = fmt.Sprintf("%s/%s/%s?page=%d&&page_size=%d", repoUrl, repoName, artifactUrl, 1, r.ArtifactCount)
		} else {
			tagsUrl = fmt.Sprintf("%s/%s/%s/%s", v1Url, repositoryUrl, r.Name, tagUrl)
		}

		body, res, err2 := c.HttpClient.Get(tagsUrl)
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
		if c.Version == "v2" {
			var artifacts []artifacts
			if err = json.Unmarshal(body, &artifacts); err != nil {
				return
			}
			for _, art := range artifacts {
				for _, tag := range art.Tags {
					items = append(items, r.Name+":"+tag.Name)
				}
			}
		} else {
			var tags []NameResult
			if err = json.Unmarshal(body, &tags); err != nil {
				return
			}
			for _, tag := range tags {
				items = append(items, r.Name+":"+tag.Name)
			}
		}
	}
	artifactEnd := artifactStart + request.Limit
	if artifactStart+request.Limit > len(items) {
		artifactEnd = len(items)
	}
	if artifactStart > artifactEnd {
		response.Items = []string{}
		return
	}
	response.Items = items[artifactStart:artifactEnd]
	if len(response.Items) == request.Limit {
		response.ContinueToken = "continue"
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
