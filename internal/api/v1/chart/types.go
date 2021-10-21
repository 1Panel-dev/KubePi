package chart

import (
	v1Chart "github.com/KubeOperator/kubepi/internal/model/v1/chart"
)

type Repo struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type RepoCreate struct {
	v1Chart.RepoCreate
}

type RepoUpdate struct {
	v1Chart.RepoUpdate
}

type Chart struct {
	Name        string `json:"name"`
	Repo        string `json:"repo"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

type Detail struct {
	Name     string `json:"name"`
	Versions string `json:"versions"`
}

type ChInstall struct {
	Name         string                 `json:"name"`
	ChartName    string                 `json:"chartName"`
	Repo         string                 `json:"repo"`
	ChartVersion string                 `json:"chartVersion"`
	Cluster      string                 `json:"cluster"`
	Values       map[string]interface{} `json:"values"`
	Namespace    string                 `json:"namespace"`
}

type HelmInstalled struct {
}
