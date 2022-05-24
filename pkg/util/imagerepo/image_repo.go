package imagerepo

import "github.com/KubeOperator/kubepi/pkg/util/imagerepo/repos"

type RepoClient interface {
	ListRepos(request repos.ProjectRequest) ([]string, error)
	ListImages(request repos.RepoRequest) (response repos.RepoResponse, err error)
	ListImagesWithoutPage(repository string) (images []string, err error)
}

func NewClient(config repos.Config) RepoClient {
	switch config.Type {
	case "Nexus":
		return repos.NewNexusClient(config.EndPoint, config.Credential.Username, config.Credential.Password)
	case "Harbor":
		return repos.NewHarborClient(config.EndPoint, config.Credential.Username, config.Credential.Password, config.Version)
	case "DockerRegistry":
		return repos.NewDockerRegistryClient(config.EndPoint, config.Credential.Username, config.Credential.Password)
	}
	return nil
}
