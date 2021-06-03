package cluster

import (
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	"github.com/KubeOperator/ekko/internal/server"
)

type Cluster interface {
	Create(cluster *v1Cluster.Cluster) error
	Get(name string) (*v1Cluster.Cluster, error)
	All() ([]v1Cluster.Cluster, error)
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
