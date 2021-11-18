package imagerepo

import "github.com/KubeOperator/kubepi/pkg/util/imagerepo/repos"

type RepoClient interface {
	//Auth() error
	ListRepos() ([]string, error)
	ListImages(repository string) (images []string, err error)
}

func NewClient(config Config) RepoClient {
	switch config.Type {
	case "Nexus":
		return repos.NewNexusClient(config.EndPoint, config.Credential.Username, config.Credential.Password)
	case "Harbor":
		return repos.NewHarborClient(config.EndPoint, config.Credential.Username, config.Credential.Password)
	case "DockerRegistry":
		return repos.NewDockerRegistryClient(config.EndPoint, config.Credential.Username, config.Credential.Password)
	}
	return nil
}

type Config struct {
	Type     string
	EndPoint string
	//DownloadUrl string
	Credential Credential
}

type Credential struct {
	Username string
	Password string
}
