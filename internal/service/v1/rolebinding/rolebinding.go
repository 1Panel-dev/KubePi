package rolebinding

import (
	"errors"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	customStorm "github.com/KubeOperator/ekko/pkg/storm"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	GetRoleBindingBySubject(subject v1Role.Subject, options common.DBOptions) ([]v1Role.Binding, error)
	CreateRoleBinding(binding *v1Role.Binding, options common.DBOptions) error
	Delete(name string, options common.DBOptions) error
}

func NewService() Service {
	return &service{
	}
}

type service struct {
	common.DefaultDBService
}

func (s *service) CreateRoleBinding(binding *v1Role.Binding, options common.DBOptions) error {
	db := s.GetDB(options)
	binding.UUID = uuid.New().String()
	binding.CreateAt = time.Now()
	binding.UpdateAt = time.Now()
	return db.Save(binding)
}

func (s *service) GetRoleBindingBySubject(subject v1Role.Subject, options common.DBOptions) ([]v1Role.Binding, error) {
	db := s.GetDB(options)
	query := db.Select(customStorm.Containers("Subjects", subject))
	var rbs []v1Role.Binding
	if err := query.Find(&rbs); err != nil {
		return rbs, err
	}
	return rbs, nil
}

func (s *service) Delete(name string, options common.DBOptions) error {
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
