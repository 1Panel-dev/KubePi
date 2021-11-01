package ldap

import (
	v1Ldap "github.com/KubeOperator/kubepi/internal/model/v1/ldap"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	ldapClient "github.com/KubeOperator/kubepi/pkg/util/ldap"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	common.DBService
	Create(ldap *v1Ldap.Ldap, options common.DBOptions) error
	List(options common.DBOptions) ([]v1Ldap.Ldap, error)
	Update(id string, ldap *v1Ldap.Ldap, options common.DBOptions) error
	GetById(id string, options common.DBOptions) (*v1Ldap.Ldap, error)
	Delete(id string, options common.DBOptions) error
}

func NewService() Service {
	return &service{

	}
}

type service struct {
	common.DefaultDBService
}

func (l *service) Create(ldap *v1Ldap.Ldap, options common.DBOptions) error {

	lc := ldapClient.NewLdapClient(ldap.Address, ldap.Port, ldap.Username, ldap.Password)
	err := lc.Connect()
	if err != nil {
		return err
	}
	db := l.GetDB(options)
	ldap.UUID = uuid.New().String()
	ldap.CreateAt = time.Now()
	ldap.UpdateAt = time.Now()
	return db.Save(ldap)
}

func (l *service) List(options common.DBOptions) ([]v1Ldap.Ldap, error) {
	db := l.GetDB(options)
	ldap := make([]v1Ldap.Ldap, 0)
	if err := db.All(&ldap); err != nil {
		return nil, err
	}
	return ldap, nil
}

func (l *service) Update(id string, ldap *v1Ldap.Ldap, options common.DBOptions) error {
	lc := ldapClient.NewLdapClient(ldap.Address, ldap.Port, ldap.Username, ldap.Password)
	if err := lc.Connect(); err != nil {
		return err
	}
	old, err := l.GetById(id, options)
	if err != nil {
		return err
	}
	ldap.UUID = old.UUID
	ldap.CreateAt = old.CreateAt
	ldap.UpdateAt = time.Now()
	db := l.GetDB(options)
	return db.Update(ldap)
}

func (l *service) GetById(id string, options common.DBOptions) (*v1Ldap.Ldap, error) {
	db := l.GetDB(options)
	var ldap v1Ldap.Ldap
	query := db.Select(q.Eq("UUID", id))
	if err := query.First(&ldap); err != nil {
		return nil, err
	}
	return &ldap, nil
}

func (l *service) Delete(id string, options common.DBOptions) error {
	db := l.GetDB(options)
	ldap, err := l.GetById(id, options)
	if err != nil {
		return err
	}
	return db.DeleteStruct(ldap)
}
