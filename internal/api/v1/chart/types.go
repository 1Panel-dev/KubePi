package chart

import (
	v1Chart "github.com/KubeOperator/kubepi/internal/model/v1/chart"
)

type Chart struct {
	v1Chart.Chart
}

type Repo struct {
	Url  string `json:"url"`
	Name string `json:"name"`
}

type RepoCreate struct {
	v1Chart.RepoCreate
}
