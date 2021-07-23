package project

import (
	"errors"
	v1Project "github.com/KubeOperator/ekko/internal/model/v1/project"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(clusterName string, r *v1Project.Project, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1Project.Project, error)
	List(clusterName string, options common.DBOptions) ([]v1Project.Project, error)
	Delete(name string, options common.DBOptions) error
	Update(name string, role *v1Project.Project, options common.DBOptions) error
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s service) Create(clusterName string, p *v1Project.Project, options common.DBOptions) error {
	db := s.GetDB(options)
	p.UUID = uuid.New().String()
	p.CreateAt = time.Now()
	p.UpdateAt = time.Now()
	p.ClusterRef = clusterName
	return db.Save(p)
}

func (s service) Get(name string, options common.DBOptions) (*v1Project.Project, error) {
	db := s.GetDB(options)
	var p v1Project.Project
	if err := db.One("Name", name, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s service) List(clusterName string, options common.DBOptions) ([]v1Project.Project, error) {
	db := s.GetDB(options)

	rs := make([]v1Project.Project, 0)
	if err := db.Find("ClusterRef", clusterName, &rs); err != nil {
		return rs, err
	}
	return rs, nil
}

func (s service) Delete(name string, options common.DBOptions) error {
	db := s.GetDB(options)

	item, err := s.Get(name, options)
	if err != nil {
		return err
	}
	if item.BuiltIn {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(item)
}

func (s service) Update(name string, role *v1Project.Project, options common.DBOptions) error {
	panic("implement me")
}
