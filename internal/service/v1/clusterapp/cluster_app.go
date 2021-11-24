package clusterapp

import (
	v1ClusterApp "github.com/KubeOperator/kubepi/internal/model/v1/clusterapp"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(clusterApp *v1ClusterApp.ClusterApp, options common.DBOptions) error
	Get(name, cluster string, options common.DBOptions) (*v1ClusterApp.ClusterApp, error)
	Delete(name, cluster string, options common.DBOptions) error
	DeleteByCluster(cluster string, options common.DBOptions) error
}

func NewService() Service {
	return &ClusterApp{
		DefaultDBService: common.DefaultDBService{},
	}
}

type ClusterApp struct {
	common.DefaultDBService
}

func (c *ClusterApp) Create(clusterApp *v1ClusterApp.ClusterApp, options common.DBOptions) error {
	db := c.GetDB(options)
	clusterApp.UUID = uuid.New().String()
	clusterApp.CreateAt = time.Now()
	clusterApp.UpdateAt = time.Now()
	return db.Save(clusterApp)
}

func (c *ClusterApp) Get(name, cluster string, options common.DBOptions) (*v1ClusterApp.ClusterApp, error) {
	db := c.GetDB(options)
	var clusterApp v1ClusterApp.ClusterApp
	query := db.Select(q.And(q.Eq("AppName", name), q.Eq("Cluster", cluster)))
	if err := query.First(&clusterApp); err != nil {
		return nil, err
	}
	return &clusterApp, nil
}

func (c *ClusterApp) Delete(name, cluster string, options common.DBOptions) error {
	db := c.GetDB(options)
	clusterApp, err := c.Get(name, cluster, options)
	if err != nil {
		return err
	}
	return db.DeleteStruct(clusterApp)
}

func (c *ClusterApp) DeleteByCluster(cluster string, options common.DBOptions) error {
	db := c.GetDB(options)
	query := db.Select(q.Eq("Cluster", cluster))
	return query.Delete(new(v1ClusterApp.ClusterApp))
}
