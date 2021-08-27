package user

import (
	"errors"
	v1User "github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/KubeOperator/ekko/internal/service/v1/common"
	costomStorm "github.com/KubeOperator/ekko/pkg/storm"
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service interface {
	common.DBService
	Create(u *v1User.User, options common.DBOptions) error
	GetByNameOrEmail(el string, options common.DBOptions) (*v1User.User, error)
	List(options common.DBOptions) ([]v1User.User, error)
	Delete(name string, options common.DBOptions) error
	Search(num, size int, pattern string, options common.DBOptions) ([]v1User.User, int, error)
	Update(name string, u *v1User.User, options common.DBOptions) error
	UpdatePassword(name string, oldPassword string, newPassword string, options common.DBOptions) error
	ResetPassword(name string, newPassword string, options common.DBOptions) error
}

func NewService() Service {
	return &service{
	}
}

type service struct {
	common.DefaultDBService
}

func (u *service) ResetPassword(name string, newPassword string, options common.DBOptions) error {
	cu, err := u.GetByNameOrEmail(name, options)
	if err != nil {
		return err
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	cu.Authenticate.Password = string(bs)
	cu.UpdateAt = time.Now()
	db := u.GetDB(options)
	return db.Update(cu)
}

func (u *service) UpdatePassword(name string, oldPassword string, newPassword string, options common.DBOptions) error {
	cu, err := u.GetByNameOrEmail(name, options)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(cu.Authenticate.Password), []byte(oldPassword)); err != nil {
		return err
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	cu.Authenticate.Password = string(bs)
	cu.UpdateAt = time.Now()
	db := u.GetDB(options)
	return db.Update(cu)
}

func (u *service) Update(name string, us *v1User.User, options common.DBOptions) error {
	cu, err := u.GetByNameOrEmail(name, options)
	if err != nil {
		return err
	}
	//if cu.BuiltIn {
	//	return errors.New("can not delete this resource,because it created by system")
	//}
	db := u.GetDB(options)
	us.UUID = cu.UUID
	us.CreateAt = cu.CreateAt
	us.UpdateAt = time.Now()
	return db.Update(us)
}

func (u *service) Search(num, size int, pattern string, options common.DBOptions) ([]v1User.User, int, error) {
	db := u.GetDB(options)
	query := func() storm.Query {
		if pattern != "" {
			return db.Select(q.Or(
				costomStorm.Like("Name", pattern),
				costomStorm.Like("NickName", pattern),
				costomStorm.Like("Email", pattern),
			)).OrderBy("CreateAt")
		}
		return db.Select().OrderBy("CreateAt")
	}()

	if num != 0 && size != 0 {
		query.Limit(size).Skip((num - 1) * size)
	}
	count, err := query.Count(&v1User.User{})
	if err != nil {
		return nil, 0, err
	}
	users := make([]v1User.User, 0)
	if err := query.Find(&users); err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *service) GetByNameOrEmail(el string, options common.DBOptions) (*v1User.User, error) {
	db := u.GetDB(options)
	var us v1User.User
	query := db.Select(q.Or(q.Eq("Name", el), q.Eq("Email", el)))

	if err := query.First(&us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *service) List(options common.DBOptions) ([]v1User.User, error) {
	db := u.GetDB(options)
	users := make([]v1User.User, 0)
	if err := db.All(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *service) Delete(name string, options common.DBOptions) error {
	db := u.GetDB(options)
	item, err := u.GetByNameOrEmail(name, options)
	if err != nil {
		return err
	}
	if item.BuiltIn {
		return errors.New("can not delete this resource,because it created by system")
	}
	return db.DeleteStruct(item)
}

func (u *service) Create(us *v1User.User, options common.DBOptions) error {
	db := u.GetDB(options)
	us.UUID = uuid.New().String()
	us.CreateAt = time.Now()
	us.UpdateAt = time.Now()
	if us.Authenticate.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(us.Authenticate.Password), bcrypt.DefaultCost) //加密处理
		us.Authenticate.Password = string(hash)
	}
	return db.Save(us)
}
