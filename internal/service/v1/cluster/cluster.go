package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(cluster *v1Cluster.Cluster, options common.DBOptions) error
	Get(name string) (*v1Cluster.Cluster, error)
	All() ([]v1Cluster.Cluster, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error)
}

func NewClusterService() Service {
	return &cluster{}
}

type cluster struct {
	common.DefaultDBService
}

func (c *cluster) Create(cluster *v1Cluster.Cluster, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster.UUID = uuid.New().String()
	cluster.CreateAt = time.Now()
	cluster.UpdateAt = time.Now()
	return db.Save(cluster)
}

func (c *cluster) Get(name string) (*v1Cluster.Cluster, error) {
	db := server.DB()
	var cluster v1Cluster.Cluster
	if err := db.One("Name", name, &cluster); err != nil {
		return nil, err
	}
	return &cluster, nil
}

func (c *cluster) All() ([]v1Cluster.Cluster, error) {
	db := server.DB()
	var clusters []v1Cluster.Cluster
	if err := db.All(&clusters); err != nil {
		return nil, err
	}
	return clusters, nil
}

func (c *cluster) Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Cluster.Cluster, int, error) {
	db := c.GetDB(options)
	query := db.Select().Limit(size).Skip((num - 1) * size)
	count, err := query.Count(&v1Cluster.Cluster{})
	if err != nil {
		return nil, 0, err
	}
	clusters := make([]v1Cluster.Cluster, 0)
	if err := query.Find(&clusters); err != nil {
		return clusters, 0, err
	}
	return clusters, count, nil
}

func (c *cluster) Delete(name string, options common.DBOptions) error {
	db := c.GetDB(options)
	cluster, err := c.Get(name)
	if err != nil {
		return err
	}
	return db.DeleteStruct(cluster)
}
