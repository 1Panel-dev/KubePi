package ldap

import (
	"errors"
	"fmt"
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Ldap "github.com/KubeOperator/kubepi/internal/model/v1/ldap"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	v1User "github.com/KubeOperator/kubepi/internal/model/v1/user"
	"github.com/KubeOperator/kubepi/internal/server"
	"github.com/KubeOperator/kubepi/internal/service/v1/common"
	"github.com/KubeOperator/kubepi/internal/service/v1/rolebinding"
	"github.com/KubeOperator/kubepi/internal/service/v1/user"
	ldapClient "github.com/KubeOperator/kubepi/pkg/util/ldap"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Service interface {
	common.DBService
	Create(ldap *v1Ldap.Ldap, options common.DBOptions) error
	List(options common.DBOptions) ([]v1Ldap.Ldap, error)
	Update(id string, ldap *v1Ldap.Ldap, options common.DBOptions) error
	GetById(id string, options common.DBOptions) (*v1Ldap.Ldap, error)
	Delete(id string, options common.DBOptions) error
	Sync(id string, options common.DBOptions) error
}

func NewService() Service {
	return &service{
		userService:        user.NewService(),
		roleBindingService: rolebinding.NewService(),
	}
}

type service struct {
	common.DefaultDBService
	userService        user.Service
	roleBindingService rolebinding.Service
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

func (l *service) Sync(id string, options common.DBOptions) error {
	ldap, err := l.GetById(id, options)
	if err != nil {
		return err
	}
	lc := ldapClient.NewLdapClient(ldap.Address, ldap.Port, ldap.Username, ldap.Password)
	if err := lc.Connect(); err != nil {
		return err
	}
	go func() {
		entries, _ := lc.Search(ldap.Dn, ldap.Filter)
		for _, entry := range entries {
			us := new(v1User.User)
			for _, at := range entry.Attributes {
				if at.Name == "cn" {
					us.Name = strings.Trim(at.Values[0], " ")
				}
				if at.Name == "mail" {
					us.Email = strings.Trim(at.Values[0], " ")
				}
				if at.Name == "sAMAccountName" {
					us.NickName = strings.Trim(at.Values[0], " ")
				}
			}
			if us.Email == "" {
				continue
			}
			if us.NickName == "" {
				us.NickName = us.Name
			}
			us.Type = v1User.LDAP
			_, err := l.userService.GetByNameOrEmail(us.Name, options)
			if errors.Is(err, storm.ErrNotFound) {
				tx, err := server.DB().Begin(true)
				if err != nil {
					server.Logger().Errorf("create tx err:  %s", err)
				}
				err = l.userService.Create(us, common.DBOptions{DB: tx})
				if err != nil {
					_ = tx.Rollback()
					server.Logger().Errorf("can not insert user %s , err:  %s", us.Name, err)
					continue
				}
				roleName := "Common User"
				binding := v1Role.Binding{
					BaseModel: v1.BaseModel{
						Kind:       "RoleBind",
						ApiVersion: "v1",
						CreatedBy:  "admin",
					},
					Metadata: v1.Metadata{
						Name: fmt.Sprintf("role-binding-%s-%s", roleName, us.Name),
					},
					Subject: v1Role.Subject{
						Kind: "User",
						Name: us.Name,
					},
					RoleRef: roleName,
				}
				if err := l.roleBindingService.CreateRoleBinding(&binding, common.DBOptions{DB: tx}); err != nil {
					_ = tx.Rollback()
					server.Logger().Errorf("can not create  user role %s , err:  %s", us.Name, err)
					continue
				}
				_ = tx.Commit()
			}
		}
	}()

	return nil
}
