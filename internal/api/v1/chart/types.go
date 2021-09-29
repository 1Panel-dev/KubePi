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

type Chart struct {
	Name        string `json:"name"`
	Repo        string `json:"repo"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}
