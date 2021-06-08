package role

import (
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/server"
	pkgV1 "github.com/KubeOperator/ekko/pkg/api/v1"
	"github.com/asdine/storm/v3/q"
	"strings"
)

type Service interface {
	Create(r *v1Role.Role) error
	Get(name string) (*v1Role.Role, error)
	List() ([]v1Role.Role, error)
	Delete(name string) error
	Search(num, size int, conditions pkgV1.Conditions) ([]v1Role.Role, int, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s service) Create(r *v1Role.Role) error {
	if r.ApiVersion == "" {
		r.ApiVersion = "v1"
	}
	if r.Kind == "" {
		r.Kind = "Role"
	}
	db := server.DB()
	return db.Save(&r)
}

func (s service) Get(name string) (*v1Role.Role, error) {
	db := server.DB()
	var r v1Role.Role
	if err := db.One("Name", name, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func (s service) List() ([]v1Role.Role, error) {
	db := server.DB()
	rs := make([]v1Role.Role, 0)
	if err := db.All(&rs); err != nil {
		return rs, err
	}
	return rs, nil
}

func (s service) Delete(name string) error {
	db := server.DB()
	return db.DeleteStruct(v1Role.Role{Metadata: v1.Metadata{Name: name}})
}

func (s service) Search(num, size int, conditions pkgV1.Conditions) ([]v1Role.Role, int, error) {
	db := server.DB()
	var ms []q.Matcher
	for key := range conditions {
		if strings.ToLower(conditions[key].Operator) == "in" {
			m := q.In(conditions[key].Field, conditions[key].Value)
			ms = append(ms, m)
		}
	}
	query := db.Select(ms...)
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
