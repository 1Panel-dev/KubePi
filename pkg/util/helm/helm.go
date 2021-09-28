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
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"io/ioutil"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	helmDriver = "configmap"
)

func nolog(format string, v ...interface{}) {}

type Config struct {
	Host          string
	BearerToken   string
	OldNamespace  string
	Namespace     string
	Architectures string
}
type Client struct {
	installActionConfig   *action.Configuration
	unInstallActionConfig *action.Configuration
	Namespace             string
	settings              *cli.EnvSettings
	Architectures         string
}

func GetSettings() *cli.EnvSettings {
	return &cli.EnvSettings{
		PluginsDirectory: helmpath.DataPath("plugins"),
		RegistryConfig:   helmpath.ConfigPath("registry.json"),
		RepositoryConfig: helmpath.ConfigPath("repositories.yaml"),
		RepositoryCache:  helmpath.CachePath("repository"),
	}
}

func NewClient(config *Config) (*Client, error) {
	client := Client{
		Architectures: config.Architectures,
	}
	client.settings = GetSettings()
	cf := genericclioptions.NewConfigFlags(true)
	inscure := true
	apiServer := config.Host
	cf.APIServer = &apiServer
	cf.BearerToken = &config.BearerToken
	cf.Insecure = &inscure
	if config.Namespace == "" {
		client.Namespace = "default"
	} else {
		client.Namespace = config.Namespace
	}
	cf.Namespace = &client.Namespace
	installActionConfig := new(action.Configuration)
	if err := installActionConfig.Init(cf, client.Namespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.installActionConfig = installActionConfig
	unInstallActionConfig := new(action.Configuration)
	if err := unInstallActionConfig.Init(cf, config.OldNamespace, helmDriver, nolog); err != nil {
		return nil, err
	}
	client.unInstallActionConfig = unInstallActionConfig
	return &client, nil
}
func (c Client) Install(name, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	client := action.NewInstall(c.installActionConfig)
	client.ReleaseName = name
	client.Namespace = c.Namespace
	client.ChartPathOptions.InsecureSkipTLSverify = true
	if len(chartVersion) != 0 {
		client.ChartPathOptions.Version = chartVersion
	}
	p, err := client.ChartPathOptions.LocateChart(chartName, c.settings)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("locate chart %s failed: %v", chartName, err))
	}
	ct, err := loader.Load(p)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("load chart %s failed: %v", chartName, err))
	}

	release, err := client.Run(ct, values)
	if err != nil {
		return release, errors.Wrap(err, fmt.Sprintf("install tool %s with chart %s failed: %v", name, chartName, err))
	}
	return release, nil
}

func (c Client) Upgrade(name, chartName, chartVersion string, values map[string]interface{}) (*release.Release, error) {
	client := action.NewUpgrade(c.installActionConfig)
	client.Namespace = c.Namespace
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
	client := action.NewUninstall(c.unInstallActionConfig)
	res, err := client.Run(name)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("uninstall tool %s failed: %v", name, err))
	}
	return res, nil
}

func (c Client) ListCharts(pattern string) ([]*search.Result, error) {
	repos, err := c.ListRepo()
	if err != nil {
		return nil, fmt.Errorf("list chart failed: %v", err)
	}
	i := search.NewIndex()
	for _, re := range repos {
		settings := GetSettings()
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

func (c Client) ListRepo() ([]*repo.Entry, error) {
	settings := GetSettings()
	var repos []*repo.Entry
	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return repos, err
	}
	return f.Repositories, nil
}

func (c Client) AddRepo(name string, url string, username string, password string) error {
	settings := GetSettings()

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
