package role

import (
	"errors"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Service interface {
	common.DBService
	Create(r *v1Role.Role, options common.DBOptions) error
	CreateWithTemplate(r *v1Role.Role, templateName string, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1Role.Role, error)
	List(options common.DBOptions) ([]v1Role.Role, error)
	Delete(name string, options common.DBOptions) error
	Update(name string, role *v1Role.Role, options common.DBOptions) error
	Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Role.Role, int, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
	common.DefaultDBService
}

func (s *service) Update(name string, role *v1Role.Role, options common.DBOptions) error {
	db := s.GetDB(options)
	r, err := s.Get(name, options)
	if err != nil {
		return err
	}
	if r.BuiltIn {
		return errors.New("can not delete this resource,because it created by system")
	}
	role.UUID = r.UUID
	role.CreateAt = r.CreateAt
	role.UpdateAt = time.Now()
	return db.Update(role)
}

func (s *service) Create(r *v1Role.Role, options common.DBOptions) error {
	db := s.GetDB(options)
	r.UUID = uuid.New().String()
	r.CreateAt = time.Now()
	r.UpdateAt = time.Now()
	return db.Save(r)
}

func (s *service) CreateWithTemplate(r *v1Role.Role, templateName string, options common.DBOptions) error {
	db := s.GetDB(options)
	var template v1Role.Role
	if err := db.One("Name", templateName, &template); err != nil {
		return err
	}
	r.Rules = template.Rules
	return s.Create(r, options)
}

func (s *service) Get(name string, options common.DBOptions) (*v1Role.Role, error) {
	db := s.GetDB(options)
	var r v1Role.Role
	if err := db.One("Name", name, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s *service) List(options common.DBOptions) ([]v1Role.Role, error) {
	db := s.GetDB(options)

	rs := make([]v1Role.Role, 0)
	if err := db.All(&rs); err != nil {
		return rs, err
	}
	return rs, nil
}

func (s *service) Delete(name string, options common.DBOptions) error {
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

func (s *service) Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Role.Role, int, error) {
	db := s.GetDB(options)
	var ms []q.Matcher
	for key := range conditions {
		if strings.ToLower(conditions[key].Operator) == "in" {
			m := q.In(conditions[key].Field, conditions[key].Value)
			ms = append(ms, m)
		}
	}
	query := db.Select(ms...).OrderBy( "CreateAt")
	if num != 0 && size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}

	count, err := query.Count(&v1Role.Role{})
	if err != nil {
		return nil, 0, err
	}
	roles := make([]v1Role.Role, 0)
	if err := query.Find(&roles); err != nil {
		return nil, 0, err
	}
	return roles, count, nil
}
