package user

import (
	v1User "github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/KubeOperator/ekko/internal/server"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"os/user"
)

type Service interface {
	Create(u *v1User.User) error
	Get(name string) (*v1User.User, error)
	GetByEmail(email string) (*v1User.User, error)
	List() ([]v1User.User, error)
	Delete(name string) error
	Search(num, size int, conditions pkgV1.Conditions) ([]user.User, int, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (u *service) Search(num, size int, conditions pkgV1.Conditions) ([]user.User, int, error) {
	db := server.DB()
	query := db.Select().Limit(size).Skip((num - 1) * size)
	count, err := query.Count(&user.User{})
	if err != nil {
		return nil, 0, err
	}
	var users []user.User
	if err := query.Find(&users); err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (u *service) Get(name string) (*v1User.User, error) {
	db := server.DB()
	var us v1User.User
	if err := db.One("Name", name, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *service) GetByEmail(email string) (*v1User.User, error) {
	db := server.DB()
	var us v1User.User
	if err := db.One("Email", email, &us); err != nil {
		return nil, err
	}
	return &us, nil
}

func (u *service) List() ([]v1User.User, error) {
	db := server.DB()
	var users []v1User.User
	if err := db.All(users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *service) Delete(name string) error {
	db := server.DB()
	return db.DeleteStruct(user.User{Name: name})
}

func (u *service) Create(us *v1User.User) error {
	db := server.DB()
	return db.Save(&us)
}
