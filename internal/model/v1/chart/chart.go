package chart

import "helm.sh/helm/v3/pkg/chart"

type RepoCreate struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ChArrayResult struct {
	Versions []string    `json:"versions"`
	Chart    chart.Chart `json:"chart"`
}
