package rolebinding

import (
	"errors"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	GetRoleBindingBySubject(subject v1Role.Subject, options common.DBOptions) ([]v1Role.Binding, error)
	GetRoleBindingsByRoleName(roleName string, options common.DBOptions) ([]v1Role.Binding, error)
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

func (s *service) GetRoleBindingsByRoleName(roleName string, options common.DBOptions) ([]v1Role.Binding, error) {
	db := s.GetDB(options)
	var r []v1Role.Binding
	if err := db.Find("RoleRef", roleName, &r); err != nil {
		return nil, err
	}
	return r, nil
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
	query := db.Select(q.Eq("Subject", subject))
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
	if binding.BuiltIn {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(&binding)
}
