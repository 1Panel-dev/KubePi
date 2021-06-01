package user

import (
	"github.com/KubeOperator/ekko/pkg/model/v1/user"
	"github.com/KubeOperator/ekko/pkg/server"
)

type Service interface {
	Create(u *user.User) error
	//Get(name string) (*user.User, error)
	//List(name string) ([]*user.User, error)
	//Update(name string, u *user.User) error
	//Delete(name string) error
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (u *service) Create(us *user.User) error {
	db := server.DB()
	return db.Save(&us)
}
