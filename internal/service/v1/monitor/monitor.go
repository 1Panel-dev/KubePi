package monitor

import (
	v1Monitor "github.com/KubeOperator/kubepi/internal/model/v1/monitor"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	costomStorm "github.com/KubeOperator/kubepi/pkg/storm"
	grafanaClient "github.com/KubeOperator/kubepi/pkg/util/grafana"
	grafana "github.com/KubeOperator/kubepi/pkg/util/grafana/dashboards"
	"github.com/KubeOperator/kubepi/pkg/util/lang"
	"github.com/KubeOperator/kubepi/pkg/util/prometheus"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	// 仪表盘
	GrafanaTestConnect(monitor *v1Monitor.GrafanaConfig) error                                 // grafana连通性测试
	GrafanaList(options common.DBOptions) (*v1Monitor.GrafanaConfig, error)                    // 查看grafana配置
	GrafanaCreate(monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error            // 创建grafana配置
	GrafanaUpdate(id string, monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error // 更新grafana配置
	GrafanaImportDashboards(monitor *v1Monitor.GrafanaConfig) error                            // grafana导入默认仪表盘

	// 指标
	MetricsSearch(num, size int, conditions common.Conditions, options common.DBOptions) (result []v1Monitor.MetricsConfig, count int, err error) // 搜索metrics实例
	MetricsCreate(metr *v1Monitor.MetricsConfig, options common.DBOptions) error                                                                  // 创建Metrics实例
	MetricsDelete(name string, options common.DBOptions) error                                                                                    // 删除Metrics实例
	MetricsUpdate(name string, metr *v1Monitor.MetricsConfig, options common.DBOptions) error                                                     // 更新Metrics实例
	MetricsGetByName(name string, options common.DBOptions) (metr v1Monitor.MetricsConfig, err error)                                             // 根据名称获取Metrics实例
	MetricsExplorer(name string, options common.DBOptions) (data interface{}, err error)                                                          // 根据Metrics实例名称，到处所有Metrics
	MetricsTestConnect(name string, options common.DBOptions) error                                                                               // prometheus连通性测试
	MetricsQuery(name, promql, timestamp string, options common.DBOptions) (data interface{}, err error)                                          // prometheus promql查询（默认为即时查询）
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s *service) GrafanaTestConnect(monitor *v1Monitor.GrafanaConfig) error {
	gc := grafanaClient.NewGrafanaClient(monitor.Address, monitor.ServiceAccountToken, monitor.Enable, monitor.DefaultDashboard)
	if err := gc.TestConnect(monitor.Address); err != nil {
		return err
	}
	return nil
}

func (s *service) GrafanaStatus(options common.DBOptions) bool {
	db := s.GetDB(options)
	monitor := make([]v1Monitor.GrafanaConfig, 0)
	if err := db.All(&monitor); err != nil {
		return false
	}

	if len(monitor) == 0 {
		return false
	}

	return monitor[0].Enable
}

func (s *service) GrafanaDashboardStatus(options common.DBOptions) bool {
	db := s.GetDB(options)
	monitor := make([]v1Monitor.GrafanaConfig, 0)
	if err := db.All(&monitor); err != nil {
		return false
	}

	if len(monitor) == 0 {
		return false
	}

	return monitor[0].DefaultDashboard
}

func (s *service) GrafanaList(options common.DBOptions) (*v1Monitor.GrafanaConfig, error) {
	db := s.GetDB(options)
	monitor := make([]v1Monitor.GrafanaConfig, 0)
	if err := db.All(&monitor); err != nil {
		return nil, err
	}

	if len(monitor) == 0 {
		var monitorNoData v1Monitor.GrafanaConfig
		return &monitorNoData, nil
	}

	return &monitor[0], nil
}

func (s *service) GrafanaCreate(monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error {
	gc := grafanaClient.NewGrafanaClient(monitor.Address, monitor.ServiceAccountToken, monitor.Enable, monitor.DefaultDashboard)
	// 当用户进行Grafana配置时，应该为用户检测目标是否可连接
	if err := gc.TestConnect(monitor.Address); err != nil {
		return err
	}

	db := s.GetDB(options)
	monitor.UUID = uuid.New().String()
	monitor.CreateAt = time.Now()
	monitor.UpdateAt = time.Now()
	return db.Save(monitor)
}

func (s *service) GrafanaUpdate(id string, monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error {
	gc := grafanaClient.NewGrafanaClient(monitor.Address, monitor.ServiceAccountToken, monitor.Enable, monitor.DefaultDashboard)
	// 当用户进行Grafana配置更新时，应该当功能为开启时才检测连接是否能通信
	if monitor.Enable {
		if err := gc.TestConnect(monitor.Address); err != nil {
			return err
		}
	}

	old, err := s.GrafanaGetById(id, options)
	if err != nil {
		return err
	}
	monitor.UUID = old.UUID
	monitor.CreateAt = old.CreateAt
	monitor.UpdateAt = time.Now()
	db := s.GetDB(options)
	if monitor.Enable != old.Enable {
		err = db.UpdateField(monitor, "Enable", monitor.Enable)
		if err != nil {
			return err
		}
	}
	return db.Update(monitor)
}

func (s *service) GrafanaGetById(id string, options common.DBOptions) (*v1Monitor.GrafanaConfig, error) {
	db := s.GetDB(options)
	var monitor v1Monitor.GrafanaConfig
	query := db.Select(q.Eq("UUID", id))
	if err := query.First(&monitor); err != nil {
		return nil, err
	}
	return &monitor, nil
}

func (s *service) GrafanaImportDashboards(monitor *v1Monitor.GrafanaConfig) error {
	gc := grafanaClient.NewGrafanaClient(monitor.Address, monitor.ServiceAccountToken, monitor.Enable, monitor.DefaultDashboard)
	if err := gc.TestConnect(monitor.Address); err != nil {
		return err
	}

	// 创建文件夹
	if err := gc.CreateFolder(gc.Address, gc.ServiceAccountToken, "KubePi Dashboards"); err != nil {
		return err
	}

	for _, dashboard := range v1Monitor.GrafanaDashboardUid {
		switch dashboard {
		case "NamespaceOverviewKubePi":
			if err := gc.ImportDashboards(gc.Address, gc.ServiceAccountToken, "k8s_namesapce_overview_kubepi", grafana.NamespaceOverviewKubePi); err != nil {
				return err
			}
		case "PodsOverviewKubePi":
			if err := gc.ImportDashboards(gc.Address, gc.ServiceAccountToken, "k8s_pod_overview_kubepi", grafana.PodsOverviewKubePi); err != nil {
				return err
			}
		}
		time.Sleep(500 * time.Millisecond) // 视乎通过API导入仪表盘需要点时间，因此每一次成功导入都等待500ms
	}
	return nil
}

func (s *service) MetricsSearch(num, size int, conditions common.Conditions, options common.DBOptions) (result []v1Monitor.MetricsConfig, count int, err error) {
	db := s.GetDB(options)
	var ms []q.Matcher
	for k := range conditions {
		if conditions[k].Field == "quick" {
			ms = append(ms, q.Or(
				costomStorm.Like("Name", conditions[k].Value),
			))
		} else {
			field := lang.FirstToUpper(conditions[k].Field)
			value := lang.ParseValueType(conditions[k].Value)

			switch conditions[k].Operator {
			case "eq":
				ms = append(ms, q.Eq(field, value))
			case "ne":
				ms = append(ms, q.Not(q.Eq(field, value)))
			case "like":
				ms = append(ms, costomStorm.Like(field, value.(string)))
			case "not like":
				ms = append(ms, q.Not(costomStorm.Like(field, value.(string))))
			}
		}
	}
	query := db.Select(ms...).OrderBy("CreateAt").Reverse()
	count, err = query.Count(&v1Monitor.MetricsConfig{})
	if err != nil {
		return
	}
	if size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	if err = query.Find(&result); err != nil {
		return
	}
	return
}

func (s *service) MetricsCreate(metr *v1Monitor.MetricsConfig, options common.DBOptions) error {
	db := s.GetDB(options)
	metr.UUID = uuid.New().String()
	metr.CreateAt = time.Now()
	metr.UpdateAt = time.Now()
	return db.Save(metr)
}

func (s *service) MetricsDelete(name string, options common.DBOptions) error {
	db := s.GetDB(options)
	item, err1 := s.MetricsGetByName(name, options)
	if err1 != nil {
		return err1
	}
	return db.DeleteStruct(&item)
}

func (s *service) MetricsGetByName(name string, options common.DBOptions) (metr v1Monitor.MetricsConfig, err error) {
	db := s.GetDB(options)
	query := db.Select(q.Eq("Name", name))
	if err = query.First(&metr); err != nil {
		return
	}
	return
}

func (s *service) MetricsUpdate(name string, metr *v1Monitor.MetricsConfig, options common.DBOptions) error {
	db := s.GetDB(options)
	old, err1 := s.MetricsGetByName(name, options)
	if err1 != nil {
		return err1
	}
	metr.UUID = old.UUID
	metr.CreateAt = old.CreateAt
	metr.UpdateAt = time.Now()

	if !old.Auth {
		metr.Credential.Password = ""
		metr.Credential.Username = ""
		metr.Credential = v1Monitor.Credential{}
		if err := db.UpdateField(metr, "Credential", metr.Credential); err != nil {
			return err
		}
	}

	if old.Auth != metr.Auth {
		if err := db.UpdateField(metr, "Auth", metr.Auth); err != nil {
			return err
		}
	}

	return db.Update(metr)
}

func (s *service) MetricsExplorer(name string, options common.DBOptions) (data interface{}, err error) {
	metr, err := s.MetricsGetByName(name, options)
	if err != nil {
		return nil, err
	}

	mc := prometheus.NewPrometheusClient(metr.EndPoint, metr.Credential.Username, metr.Credential.Password)
	return mc.GetMetrics(mc.Address, mc.Username, mc.Username)
}

func (s *service) MetricsTestConnect(name string, options common.DBOptions) error {
	metr, err := s.MetricsGetByName(name, options)
	if err != nil {
		return err
	}

	mc := prometheus.NewPrometheusClient(metr.EndPoint, metr.Credential.Username, metr.Credential.Password)
	return mc.TestConnect(mc.Address, mc.Username, mc.Password)
}

func (s *service) MetricsQuery(name, promql, timestamp string, options common.DBOptions) (data interface{}, err error) {
	metr, err := s.MetricsGetByName(name, options)
	if err != nil {
		return nil, err
	}

	mc := prometheus.NewPrometheusClient(metr.EndPoint, metr.Credential.Username, metr.Credential.Password)
	return mc.QueryMetrics(mc.Address, mc.Username, mc.Password, promql, timestamp)
}
