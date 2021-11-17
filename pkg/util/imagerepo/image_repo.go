package imagerepo

import "github.com/KubeOperator/kubepi/pkg/util/imagerepo/repos"

type RepoClient interface {
	Auth() error
	ListRepos() ([]string, error)
}

func NewClient(config Config) RepoClient {
	switch config.Type {
	case "Nexus":
		return repos.NewNexusClient(config.EndPoint, config.Credential.Username, config.Credential.Password, config.DownloadUrl)
	}

	return nil
}

type Config struct {
	Type        string
	EndPoint    string
	DownloadUrl string
	Credential  Credential
}

type Credential struct {
	Username string
	Password string
}
