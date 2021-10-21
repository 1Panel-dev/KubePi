package chart

import (
	"helm.sh/helm/v3/pkg/chart"
	"time"
)

type RepoCreate struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type ChArrayResult struct {
	Versions []Version `json:"versions"`
	Chart    ChDetail  `json:"chart"`
}

type ChDetail struct {
	Metadata chart.Metadata         `json:"metadata"`
	Readme   string                 `json:"readme"`
	Values   map[string]interface{} `json:"values"`
}

type Version struct {
	Version string    `json:"version"`
	Date    time.Time `json:"date"`
}

type UpdateVersion struct {
	Version    string `json:"version"`
	AppVersion string `json:"appVersion"`
}

type UpdateResult struct {
	Repo     string          `json:"repo"`
	Versions []UpdateVersion `json:"versions"`
}

type RepoUpdate struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Repo struct {
	Url      string `json:"url"`
	Name     string `json:"name"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}
