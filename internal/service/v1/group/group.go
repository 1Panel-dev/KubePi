package group

import (
	"errors"
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(u *v1Group.Group, options common.DBOptions) error
	Update(name string, u *v1Group.Group, options common.DBOptions) error
	Get(name string, options common.DBOptions) (*v1Group.Group, error)
	List(options common.DBOptions) ([]v1Group.Group, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Group.Group, int, error)
}

func NewService() Service {
	return &service{
		DefaultDBService: common.DefaultDBService{},
	}
}

type service struct {
	common.DefaultDBService
}


func (s *service) Update(name string, g *v1Group.Group, options common.DBOptions) error {
	cu, err := s.Get(name, options)
	if err != nil {
		return err
	}
	if cu.CreatedBy == "system" {
		return errors.New("can not delete this resource,because it created by system")
	}
	db := s.GetDB(options)
	g.UUID = cu.UUID
	g.CreateAt = cu.CreateAt
	g.UpdateAt = time.Now()
	return db.Save(g)
}

func (s *service) Get(name string, options common.DBOptions) (*v1Group.Group, error) {
	db := s.GetDB(options)
	var g v1Group.Group
	if err := db.One("Name", name, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *service) Delete(name string, options common.DBOptions) error {
	db := s.GetDB(options)
	item, err := s.Get(name, options)
	if err != nil {
		return err
	}
	if item.CreatedBy == "system" {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(item)
}

func (s *service) Create(g *v1Group.Group, options common.DBOptions) error {
	db := s.GetDB(options)
	g.UUID = uuid.New().String()
	g.CreateAt = time.Now()
	g.UpdateAt = time.Now()
	return db.Save(g)
}

func (s *service) List(options common.DBOptions) ([]v1Group.Group, error) {
	db := s.GetDB(options)
	gs := make([]v1Group.Group, 0)
	if err := db.All(&gs); err != nil {
		return gs, err
	}
	return gs, nil
}

func (s *service) Search(num, size int, conditions pkgV1.Conditions, options common.DBOptions) ([]v1Group.Group, int, error) {

	db := s.GetDB(options)
	query := db.Select()
	if num != 0 && size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	count, err := query.Count(&v1Group.Group{})
	if err != nil {
		return nil, 0, err
	}
	groups := make([]v1Group.Group, 0)
	if err := query.Find(&groups); err != nil {
		return nil, 0, err
	}
	return groups, count, nil

}
