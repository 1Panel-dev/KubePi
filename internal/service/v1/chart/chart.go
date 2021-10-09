package chart

import (
	"errors"
	v1Chart "github.com/KubeOperator/kubepi/internal/model/v1/chart"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/util/helm"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/repo"
)

type Service interface {
	common.DBService
	SearchRepo(cluster string) ([]*repo.Entry, error)
	AddRepo(cluster string, create *v1Chart.RepoCreate) error
	ListCharts(cluster, repo string, num, size int, pattern string) ([]*search.Result, int, error)
	RemoveRepo(cluster string, name string) error
	GetCharts(cluster, repo, name string) (*v1Chart.ChArrayResult, error)
	GetChartByVersion(cluster, repo, name, version string) (*v1Chart.ChDetail, error)
	InstallChart(cluster, repoName, name, chartName, chartVersion string, values map[string]interface{}) error
}

func NewService() Service {
	return &service{
		clusterService: cluster.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	clusterService cluster.Service
}

func (c *service) SearchRepo(cluster string) ([]*repo.Entry, error) {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return nil, err
	}
	repos, err := helmClient.ListRepo()
	if err != nil {
		return nil, err
	}
	return repos, err
}

func (c *service) AddRepo(cluster string, create *v1Chart.RepoCreate) error {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return err
	}
	err = helmClient.AddRepo(create.Name, create.Url, create.UserName, create.Password)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) RemoveRepo(cluster string, name string) error {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return err
	}
	success, err := helmClient.RemoveRepo(name)
	if err != nil {
		return err
	}
	if !success {
		return errors.New("delete repo failed!")
	}
	return nil
}

func (c *service) ListCharts(cluster, repo string, num, size int, pattern string) ([]*search.Result, int, error) {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return nil, 0, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return nil, 0, err
	}
	charts, err := helmClient.ListCharts(repo, pattern)
	if err != nil {
		return nil, 0, err
	}
	var chartArray []*search.Result
	for _, chart := range charts {
		exist := false
		for _, ca := range chartArray {
			if ca.Name == chart.Name {
				exist = true
				break
			}
		}
		if exist {
			continue
		}
		chartArray = append(chartArray, chart)
	}
	end := num * size
	if end > len(chartArray) {
		end = len(chartArray)
	}
	result := chartArray[(num-1)*size : end]
	return result, len(chartArray), nil
}

func (c *service) GetCharts(cluster, repo, name string) (*v1Chart.ChArrayResult, error) {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return nil, err
	}
	allVersionCharts, err := helmClient.GetCharts(repo, name)
	if err != nil {
		return nil, err
	}
	var result v1Chart.ChArrayResult
	for _, chart := range allVersionCharts {
		result.Versions = append(result.Versions, v1Chart.Version{
			Version: chart.Chart.Version,
			Date:    chart.Chart.Created,
		})
	}
	lastVersion := allVersionCharts[0].Chart.Metadata.Version
	chart, err := helmClient.GetChartDetail(repo, allVersionCharts[0].Chart.Name, lastVersion)
	if err != nil {
		return nil, err
	}
	result.Chart.Metadata = *chart.Metadata
	result.Chart.Values = chart.Values
	for _, file := range chart.Files {
		if file.Name == "README.md" {
			result.Chart.Readme = string(file.Data)
		}
	}
	return &result, nil
}

func (c *service) GetChartByVersion(cluster, repo, name, version string) (*v1Chart.ChDetail, error) {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return nil, err
	}
	chart, err := helmClient.GetChartDetail(repo, name, version)
	if err != nil {
		return nil, err
	}
	var result v1Chart.ChDetail
	result.Values = chart.Values
	result.Metadata = *chart.Metadata
	for _, file := range chart.Files {
		if file.Name == "README.md" {
			result.Readme = string(file.Data)
		}
	}
	return &result, nil
}

func (c *service) InstallChart(cluster, repoName, name, chartName, chartVersion string, values map[string]interface{}) error {
	clu, err := c.clusterService.Get(cluster, common.DBOptions{})
	if err != nil {
		return err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		BearerToken: clu.Spec.Authentication.BearerToken,
	})
	if err != nil {
		return err
	}
	_, err = helmClient.Install(name, repoName, chartName, chartVersion, values)
	if err != nil {
		return err
	}
	return nil
}
