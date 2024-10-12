package monitor

import (
	v1Monitor "github.com/KubeOperator/kubepi/internal/model/v1/monitor"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	grafanaClient "github.com/KubeOperator/kubepi/pkg/util/grafana"
	grafana "github.com/KubeOperator/kubepi/pkg/util/grafana/dashboards"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	GrafanaTestConnect(monitor *v1Monitor.GrafanaConfig) error                                 // grafana连通性测试
	GrafanaList(options common.DBOptions) (*v1Monitor.GrafanaConfig, error)                    // 查看grafana配置
	GrafanaCreate(monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error            // 创建grafana配置
	GrafanaUpdate(id string, monitor *v1Monitor.GrafanaConfig, options common.DBOptions) error // 更新grafana配置
	GrafanaImportDashboards(monitor *v1Monitor.GrafanaConfig) error                            // grafana导入默认仪表盘

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
	// 当用户进行Grafana配置时，应该为用户检测目标是否可连接
	if err := gc.TestConnect(monitor.Address); err != nil {
		return err
	}

	old, err := s.GetById(id, options)
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

func (s *service) GetById(id string, options common.DBOptions) (*v1Monitor.GrafanaConfig, error) {
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
