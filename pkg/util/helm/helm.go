package helm

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/helmpath"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"k8s.io/cli-runtime/pkg/genericclioptions"
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
	//if err := updateRepo(c.Architectures); err != nil {
	//	return nil, err
	//}
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
	//if err := updateRepo(c.Architectures); err != nil {
	//	return nil, err
	//}
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
	release, err := client.Run(name)
	if err != nil {
		return release, errors.Wrap(err, fmt.Sprintf("uninstall tool %s failed: %v", name, err))
	}
	return release, nil
}

func (c Client) List() ([]*release.Release, error) {
	client := action.NewList(c.unInstallActionConfig)
	client.All = true
	release, err := client.Run()
	if err != nil {
		return release, fmt.Errorf("list chart failed: %v", err)
	}
	return release, nil
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
