package clusterbinding

import (
	"errors"
	v1Cluster "github.com/KubeOperator/ekko/internal/model/v1/cluster"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	GetClusterBindingByClusterName(clusterName string, options common.DBOptions) ([]v1Cluster.Binding, error)
	CreateClusterBinding(binding *v1Cluster.Binding, options common.DBOptions) error
	Delete(name string, options common.DBOptions) error
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s service) CreateClusterBinding(binding *v1Cluster.Binding, options common.DBOptions) error {
	db := s.GetDB(options)
	binding.UUID = uuid.New().String()
	binding.CreateAt = time.Now()
	binding.UpdateAt = time.Now()
	return db.Save(binding)
}

func (s service) GetClusterBindingByClusterName(clusterName string, options common.DBOptions) ([]v1Cluster.Binding, error) {
	db := s.GetDB(options)
	query := db.Select(q.Eq("ClusterRef", clusterName))
	var rbs []v1Cluster.Binding
	if err := query.Find(&rbs); err != nil {
		return rbs, err
	}
	return rbs, nil
}

func (s service) Delete(name string, options common.DBOptions) error {
	db := s.GetDB(options)
	var binding v1Role.Binding
	if err := db.One("Name", name, &binding); err != nil {
		return err
	}
	if binding.CreatedBy == "system" {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(&binding)
}
