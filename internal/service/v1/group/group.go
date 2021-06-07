package group

import (
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Group "github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/server"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
)

type Service interface {
	Create(u *v1Group.Group) error
	Get(name string) (*v1Group.Group, error)
	List() ([]v1Group.Group, error)
	Delete(name string) error
	Search(num, size int, conditions pkgV1.Conditions) ([]v1Group.Group, int, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s service) Create(u *v1Group.Group) error {
	if u.ApiVersion == "" {
		u.ApiVersion = "v1"
	}
	if u.Kind == "" {
		u.Kind = "Group"
	}
	db := server.DB()
	return db.Save(&u)
}

func (s service) Get(name string) (*v1Group.Group, error) {
	db := server.DB()
	var g v1Group.Group
	if err := db.One("Name", name, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (s service) List() ([]v1Group.Group, error) {
	db := server.DB()
	gs := make([]v1Group.Group, 0)
	if err := db.All(&gs); err != nil {
		return gs, err
	}
	return gs, nil
}

func (s service) Delete(name string) error {
	db := server.DB()
	return db.DeleteStruct(v1Group.Group{Metadata: v1.Metadata{Name: name}})
}

func (s service) Search(num, size int, conditions pkgV1.Conditions) ([]v1Group.Group, int, error) {

	db := server.DB()
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
