package chart

import (
	"errors"
	v1Chart "github.com/KubeOperator/kubepi/internal/model/v1/chart"
	v1ClusterApp "github.com/KubeOperator/kubepi/internal/model/v1/clusterapp"
	"github.com/KubeOperator/kubepi/internal/service/v1/cluster"
	"github.com/KubeOperator/kubepi/internal/service/v1/clusterapp"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/pkg/kubernetes"
	"github.com/KubeOperator/kubepi/pkg/util/helm"
	"github.com/asdine/storm/v3"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type Service interface {
	common.DBService
	SearchRepo(cluster string) ([]*repo.Entry, error)
	AddRepo(cluster string, create *v1Chart.RepoCreate) error
	UpdateRepo(cluster string, update *v1Chart.RepoUpdate) error
	GetRepo(cluster string, name string) (*v1Chart.Repo, error)
	RemoveRepo(cluster string, name string) error
	ListCharts(cluster, repo string, num, size int, pattern string) ([]*search.Result, int, error)
	GetCharts(cluster, repo, name string) (*v1Chart.ChArrayResult, error)
	GetChartByVersion(cluster, repo, name, version string) (*v1Chart.ChDetail, error)
	InstallChart(cluster, repoName, namespace, name, chartName, chartVersion string, values map[string]interface{}) error
	ListAllInstalled(cluster, namespace string, num, size int, pattern string) ([]*release.Release, int, error)
	UnInstallChart(cluster, namespace, name string) error
	GetAppDetail(cluster string, name string) (*release.Release, error)
	GetChartsUpdate(cluster, chart, name string) (*v1Chart.UpdateResult, error)
	UpgradeChart(cluster, namespace, repoName, name, chartName, chartVersion string, values map[string]interface{}) error
	SyncRepo(cluster,name string) error
}

func NewService() Service {
	return &service{
		clusterService:    cluster.NewService(),
		clusterAppService: clusterapp.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	clusterService    cluster.Service
	clusterAppService clusterapp.Service
}

func (c *service) SearchRepo(cluster string) ([]*repo.Entry, error) {
	helmClient, err := NewHelmClient(cluster, "")
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
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return err
	}
	err = helmClient.AddRepo(create.Name, create.Url, create.UserName, create.Password)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) GetRepo(cluster string, name string) (*v1Chart.Repo, error) {
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return nil, err
	}
	re, err := helmClient.GetRepo(name)
	if err != nil {
		return nil, err
	}
	return &v1Chart.Repo{Name: re.Name, Url: re.URL, UserName: re.Username, Password: re.Password}, nil
}

func (c *service) UpdateRepo(cluster string, update *v1Chart.RepoUpdate) error {
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return err
	}
	result, err := helmClient.RemoveRepo(update.Name)
	if err != nil {
		return err
	}
	if !result {
		return errors.New("repo delete failed")
	}
	err = helmClient.AddRepo(update.Name, update.Url, update.UserName, update.Password)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) RemoveRepo(cluster string, name string) error {
	helmClient, err := NewHelmClient(cluster, "")
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
	helmClient, err := NewHelmClient(cluster, "")
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
	result := []*search.Result{}
	if len(chartArray) > 0 {
		result = chartArray[(num-1)*size : end]
	}
	return result, len(chartArray), nil
}

func (c *service) GetCharts(cluster, repo, name string) (*v1Chart.ChArrayResult, error) {
	helmClient, err := NewHelmClient(cluster, "")
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
	if len(allVersionCharts) > 0 && allVersionCharts[0].Chart.Metadata.Version != "" {
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
	}
	return &result, nil
}

func (c *service) GetChartsUpdate(cluster, chart, name string) (*v1Chart.UpdateResult, error) {
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return nil, err
	}
	clusterApp, err := c.clusterAppService.Get(name, cluster, common.DBOptions{})
	if err != nil {
		return nil, errors.New("There is no version that can be upgraded")
	}
	allVersionCharts, err := helmClient.GetCharts(clusterApp.Repo, chart)
	if err != nil {
		return nil, err
	}
	update := &v1Chart.UpdateResult{
		Repo:     clusterApp.Repo,
		Versions: []v1Chart.UpdateVersion{},
	}
	for _, chart := range allVersionCharts {
		update.Versions = append(update.Versions, v1Chart.UpdateVersion{
			Version:    chart.Chart.Version,
			AppVersion: chart.Chart.AppVersion,
		})
	}
	return update, nil
}

func (c *service) GetChartByVersion(cluster, repo, name, version string) (*v1Chart.ChDetail, error) {
	helmClient, err := NewHelmClient(cluster, "")
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

func (c *service) InstallChart(cluster, repoName, namespace, name, chartName, chartVersion string, values map[string]interface{}) error {
	helmClient, err := NewHelmClient(cluster, namespace)
	if err != nil {
		return err
	}
	_, err = helmClient.Install(name, repoName, chartName, chartVersion, values)
	if err != nil {
		return err
	}
	err = c.clusterAppService.Create(&v1ClusterApp.ClusterApp{
		AppName: name,
		Repo:    repoName,
		Cluster: cluster,
	}, common.DBOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *service) UpgradeChart(cluster, namespace, repoName, name, chartName, chartVersion string, values map[string]interface{}) error {
	helmClient, err := NewHelmClient(cluster, namespace)
	if err != nil {
		return err
	}
	_, err = helmClient.Upgrade(name, repoName, chartName, chartVersion, values)
	if err != nil {
		return err
	}
	return nil
}

func (c *service) UnInstallChart(cluster, namespace, name string) error {
	helmClient, err := NewHelmClient(cluster, namespace)
	if err != nil {
		return err
	}
	_, err = helmClient.Uninstall(name)
	if err != nil {
		return err
	}
	err = c.clusterAppService.Delete(name, cluster, common.DBOptions{})
	if err != nil && err != storm.ErrNotFound {
		return err
	}
	return nil
}

func (c *service) ListAllInstalled(cluster, namespace string, num, size int, pattern string) ([]*release.Release, int, error) {
	helmClient, err := NewHelmClient(cluster, namespace)
	if err != nil {
		return nil, 0, err
	}
	releases, total, err := helmClient.List(size, num, pattern)
	if err != nil {
		return nil, 0, err
	}
	return releases, total, nil
}

func (c *service) GetAppDetail(cluster string, name string) (*release.Release, error) {
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return nil, err
	}
	return helmClient.GetDetail(name)
}

func (c *service) SyncRepo(cluster,name string) error {
	helmClient, err := NewHelmClient(cluster, "")
	if err != nil {
		return err
	}
	return helmClient.UpdateRepo(name)
}

func NewHelmClient(clusterName, namespace string) (*helm.Client, error) {
	clu, err := cluster.NewService().Get(clusterName, common.DBOptions{})
	if err != nil {
		return nil, err
	}
	kubeConfig, err := kubernetes.NewKubernetes(clu).Config()
	if err != nil {
		return nil, err
	}
	helmClient, err := helm.NewClient(&helm.Config{
		Host:        clu.Spec.Connect.Forward.ApiServer,
		ClusterName: clusterName,
		KubeConfig:  kubeConfig,
		Namespace:   namespace,
	})
	if err != nil {
		return nil, err
	}
	return helmClient, nil
}
