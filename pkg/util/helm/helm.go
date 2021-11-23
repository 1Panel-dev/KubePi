package helm

import (
	"context"
	"fmt"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/cmd/helm/search"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	helmDriver = "secrets"
)

func nolog(format string, v ...interface{}) {}

type Config struct {
	Host          string
	BearerToken   string
	OldNamespace  string
	Namespace     string
	Architectures string
	ClusterName   string
	KubeConfig    *rest.Config
}
type Client struct {
	//installActionConfig   *action.Configuration
	//unInstallActionConfig *action.Configuration
	//listActionConfig      *action.Configuration
	//getActionConfig       *action.Configuration
	actionConfig  *action.Configuration
	Namespace     string
	settings      *cli.EnvSettings
	Architectures string
	ClusterName   string
}

func GetSettings(cluster string) *cli.EnvSettings {
	checkFiles(cluster)
	return &cli.EnvSettings{
		PluginsDirectory: helmpath.DataPath(cluster + "/plugins"),
		RegistryConfig:   helmpath.ConfigPath(cluster + "/registry.json"),
		RepositoryConfig: helmpath.ConfigPath(cluster + "/repositories.yaml"),
		RepositoryCache:  helmpath.CachePath(cluster + "/repository"),
	}
}

func checkFiles(cluster string) {

	var dataPath = helmpath.ConfigPath(cluster)
	var repositoryPath = dataPath + "/repositories.yaml"
	var registryPath = dataPath + "/registry.json"

	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		os.MkdirAll(dataPath, 0777)
		os.Chmod(dataPath, 0777)
	}
	if _, err := os.Stat(repositoryPath); os.IsNotExist(err) {
		os.Create(repositoryPath)
	}
	if _, err := os.Stat(registryPath); os.IsNotExist(err) {
		os.Create(registryPath)
	}
	return
}

func NewClient(config *Config) (*Client, error) {
	client := Client{
		Architectures: config.Architectures,
	}
	client.settings = GetSettings(config.ClusterName)
	cf := genericclioptions.NewConfigFlags(true)
	apiServer := config.Host
	cf.APIServer = &apiServer
	kubeConfig := config.KubeConfig
	cf.WrapConfigFn = func(config *rest.Config) *rest.Config {
		return kubeConfig
	}
	cf.CacheDir = nil
	if config.Namespace != "" {
		client.Namespace = config.Namespace
		cf.Namespace = &client.Namespace
	}
	client.ClusterName = config.ClusterName
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(cf, config.Namespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.actionConfig = actionConfig

	return &client, nil
}

func (c Client) List(limit, offset int, pattern string) ([]*release.Release, int, error) {
	client := action.NewList(c.actionConfig)
	if c.Namespace == "" {
		client.AllNamespaces = true
		client.All = true
	}
	client.SetStateMask()
	list, err := client.Run()
	if err != nil {
		return nil, 0, err
	}
	client.Limit = limit
	client.Offset = offset - 1
	if pattern != "" {
		client.Filter = pattern
	}
	result, err := client.Run()
	if err != nil {
		return nil, 0, err
	}

	return result, len(list), nil
}

func (c Client) GetDetail(name string) (*release.Release, error) {
	client := action.NewGet(c.actionConfig)
	result, err := client.Run(name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c Client) Install(name, repoName, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var rp *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			rp = r
		}
	}
	if rp == nil {
		return nil, errors.New("get chart detail failed, repo not found")
	}
	client := action.NewInstall(c.actionConfig)
	client.ReleaseName = name
	client.Namespace = c.Namespace
	client.RepoURL = rp.URL
	client.Username = rp.Username
	client.Password = rp.Password
	client.ChartPathOptions.InsecureSkipTLSverify = true
	if len(chartVersion) != 0 {
		client.ChartPathOptions.Version = chartVersion
	}
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, fmt.Errorf("locate chart %s failed: %v", chartName, err)
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, fmt.Errorf("load chart %s failed: %v", chartName, err)
	}
	re, err := client.Run(ct, values)
	if err != nil {
		return re, errors.Wrap(err, fmt.Sprintf("install %s with chart %s failed: %v", name, chartName, err))
	}
	return re, nil
}

func (c Client) Upgrade(name, repoName, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var rp *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			rp = r
		}
	}
	if rp == nil {
		return nil, errors.New("get chart detail failed, repo not found")
	}

	client := action.NewUpgrade(c.actionConfig)
	client.Namespace = c.Namespace
	client.RepoURL = rp.URL
	client.Username = rp.Username
	client.Password = rp.Password
	client.DryRun = false
	client.ChartPathOptions.InsecureSkipTLSverify = true
	client.ChartPathOptions.Version = chartVersion
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("locate chart %s failed: %v", chartName, err))
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("load chart %s failed: %v", chartName, err))

	}

	release, err := client.Run(name, ct, values)
	if err != nil {
		return release, errors.Wrap(err, fmt.Sprintf("upgrade tool %s with chart %s failed: %v", name, chartName, err))
	}
	return release, nil
}

func (c Client) Uninstall(name string) (*release.UninstallReleaseResponse, error) {
	client := action.NewUninstall(c.actionConfig)
	res, err := client.Run(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("uninstall tool %s failed: %v", name, err))
	}
	return res, nil
}

func (c Client) ListCharts(repoName, pattern string) ([]*search.Result, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, fmt.Errorf("list chart failed: %v", err)
	}
	i := search.NewIndex()
	for _, re := range repos {
		if repoName != "KRepoAll" {
			if repoName != re.Name {
				continue
			}
		}
		settings := GetSettings(c.ClusterName)
		path := filepath.Join(settings.RepositoryCache, helmpath.CacheIndexFile(re.Name))
		indexFile, err := repo.LoadIndexFile(path)
		if err != nil {
			return nil, fmt.Errorf("list chart failed: %v", err)
		}
		i.AddRepo(re.Name, indexFile, true)
	}
	var result []*search.Result
	if pattern != "" {
		result = i.SearchLiteral(pattern, 1)
	} else {
		result = i.All()
	}
	search.SortScore(result)
	return result, nil
}

func (c Client) GetCharts(repoName, name string) ([]*search.Result, error) {
	charts, err := c.ListCharts(repoName, name)
	if err != nil {
		return nil, err
	}
	var result []*search.Result
	for _, chart := range charts {
		if chart.Chart.Name == name {
			result = append(result, chart)
		}
	}
	return result, nil
}

func (c Client) GetChartDetail(repoName, name, version string) (*chart.Chart, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var re *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			re = r
		}
	}
	if re == nil {
		return nil, errors.New("get chart detail failed, repo not found")
	}
	client := action.NewShow(action.ShowAll)
	client.Version = version
	client.RepoURL = re.URL
	client.Username = re.Username
	client.Password = re.Password
	p, err := client.LocateChart(name, c.settings)
	if err != nil {
		return nil, err
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, err
	}
	return ct, nil
}

func (c Client) ListRepo() ([]*repo.Entry, error) {
	settings := GetSettings(c.ClusterName)
	var repos []*repo.Entry
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return repos, err
	}
	return f.Repositories, nil
}

func (c Client) RemoveRepo(name string) (bool, error) {
	settings := GetSettings(c.ClusterName)
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return false, err
	}
	success := f.Remove(name)
	if !success {
		return false, err
	}
	if err := f.WriteFile(settings.RepositoryConfig, 0644); err != nil {
		return false, err
	}
	return true, nil
}

func (c Client) AddRepo(name string, url string, username string, password string) error {
	settings := GetSettings(c.ClusterName)

	repoFile := settings.RepositoryConfig

	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func() {
			if err := fileLock.Unlock(); err != nil {
				server.Logger().Errorf("addRepo fileLock.Unlock failed, error: %s", err.Error())
			}
		}()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if f.Has(name) {
		return errors.Errorf("repository name (%s) already exists, please specify a different name", name)
	}

	e := repo.Entry{
		Name:                  name,
		URL:                   url,
		Username:              username,
		Password:              password,
		InsecureSkipTLSverify: true,
	}

	r, err := repo.NewChartRepository(&e, getter.All(settings))
	if err != nil {
		return err
	}
	r.CachePath = settings.RepositoryCache
	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
	}
	f.Update(&e)
	if err := f.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

func (c Client) updateRepo(repoName string) error {
	repos, err := c.ListRepo()
	if err != nil {
		return err
	}
	var re *repo.Entry
	for _, r := range repos {
		if r.Name == repoName {
			re = r
		}
	}
	if re == nil {
		return errors.New("repo not found")
	}

	settings := GetSettings(c.ClusterName)
	repoFile := settings.RepositoryConfig
	repoCache := settings.RepositoryCache
	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return fmt.Errorf("load file of repo %s failed: %v", repoFile, err)
	}
	f.Update(re)
	var rps []*repo.ChartRepository
	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		if repoCache != "" {
			r.CachePath = repoCache
		}
		rps = append(rps, r)
	}
	updateCharts(rps)
	return nil
}

func (c Client) GetRepo(name string) (*repo.Entry, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, err
	}
	var re *repo.Entry
	for _, r := range repos {
		if r.Name == name {
			re = r
			break
		}
	}
	return re, nil
}

func updateCharts(repos []*repo.ChartRepository) {
	fmt.Printf("Hang tight while we grab the latest from your chart repositories...")
	var wg sync.WaitGroup
	for _, re := range repos {
		wg.Add(1)
		go func(re *repo.ChartRepository) {
			defer wg.Done()
			if _, err := re.DownloadIndexFile(); err != nil {
				fmt.Printf("...Unable to get an update from the %q chart repository (%s):\n\t%s\n", re.Config.Name, re.Config.URL, err)
			} else {
				fmt.Printf("...Successfully got an update from the %q chart repository\n", re.Config.Name)
			}
		}(re)
	}
	wg.Wait()
	fmt.Printf("Update Complete. ⎈ Happy Helming!⎈ ")
}
