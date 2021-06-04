package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/server"
)

type Cluster interface {
	Create(cluster *v1Cluster.Cluster) error
	Get(name string) (*v1Cluster.Cluster, error)
	All() ([]v1Cluster.Cluster, error)
	Delete(name string) error
	Search(num, size int) ([]v1Cluster.Cluster, int, error)
}

func NewClusterService() Cluster {
	return &cluster{}
}

type cluster struct {
}

func (c *cluster) Create(cluster *v1Cluster.Cluster) error {
	db := server.DB()
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

func (c *cluster) Search(num, size int) ([]v1Cluster.Cluster, int, error) {
	db := server.DB()
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

func (c *cluster) Delete(name string) error {
	db := server.DB()
	cluster, err := c.Get(name)
	if err != nil {
		return err
	}
	return db.DeleteStruct(cluster)
}
