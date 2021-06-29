package groupbinding

import (
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(binding *v1Group.Binding, options common.DBOptions) error
	Delete(name string, options common.DBOptions) error
	ListByGroupName(groupName string, options common.DBOptions) ([]v1Group.Binding, error)
	ListByUserName(userName string, options common.DBOptions) ([]v1Group.Binding, error)
	GetGroupBindingByGroupNameAndUserName(username string, groupName string, options common.DBOptions) (*v1Group.Binding, error)
}

func NewService() Service {
	return &service{
		DefaultDBService: common.DefaultDBService{},
	}
}

type service struct {
	common.DefaultDBService
}

func (s *service) GetGroupBindingByGroupNameAndUserName(username string, groupName string, options common.DBOptions) (*v1Group.Binding, error) {
	db := s.GetDB(options)
	var gb v1Group.Binding
	query := db.Select(q.And(q.Eq("UserRef", username), q.Eq("GroupRef", groupName)))
	if err := query.First(&gb); err != nil {
		return nil, err
	}
	return &gb, nil
}

func (s *service) ListByUserName(userName string, options common.DBOptions) ([]v1Group.Binding, error) {
	db := s.GetDB(options)
	bindings := make([]v1Group.Binding, 0)
	if err := db.Find("UserRef", userName, &bindings); err != nil {
		return bindings, err
	}
	return bindings, nil
}

func (s *service) Create(binding *v1Group.Binding, options common.DBOptions) error {
	db := s.GetDB(options)
	binding.UUID = uuid.New().String()
	binding.CreateAt = time.Now()
	binding.UpdateAt = time.Now()
	return db.Save(binding)
}

func (s *service) Delete(name string, options common.DBOptions) error {
	db := s.GetDB(options)
	var binding v1Group.Binding
	//if err := db.Select(q.And(q.Eq("GroupRef", groupName), q.Eq("UserRef", userName))).First(&binding); err != nil {
	//	return err
	//}
	if err := db.One("Name", name, &binding); err != nil {
		return err
	}
	return db.DeleteStruct(&binding)
}

func (s *service) ListByGroupName(groupName string, options common.DBOptions) ([]v1Group.Binding, error) {
	db := s.GetDB(options)
	bindings := make([]v1Group.Binding, 0)
	if err := db.Find("GroupRef", groupName, &bindings); err != nil {
		return bindings, err
	}
	return bindings, nil
}
