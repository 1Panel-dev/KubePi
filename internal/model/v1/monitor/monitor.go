package monitor

import v1 "github.com/KubeOperator/kubepi/internal/model/v1"

type GrafanaConfig struct {
	v1.BaseModel        `storm:"inline"`
	v1.Metadata         `storm:"inline"`
	Address             string `json:"address"`
	Enable              bool   `json:"enable"`
	DefaultDashboard    bool   `json:"default_dashboard"`
	ServiceAccountToken string `json:"service_account_token"`
}

// Grafana默认仪表盘UID
var GrafanaDashboardUid = []string{
	"NamespaceOverviewKubePi",
	"PodsOverviewKubePi",
}

type MetricsConfig struct {
	v1.BaseModel `storm:"inline"`
	v1.Metadata  `storm:"inline"`
	Type         string     `json:"type"`
	EndPoint     string     `json:"endPoint"`
	Credential   Credential `json:"credential"`
	Auth         bool       `json:"auth"`
}

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
