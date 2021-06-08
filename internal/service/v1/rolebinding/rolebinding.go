package rolebinding

import (
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	customStorm "github.com/KubeOperator/ekko/pkg/storm"
	"github.com/KubeOperator/ekko/internal/server"
)

type Service interface {
	GetRoleBindingBySubject(subject v1Role.Subject) ([]v1Role.Banding, error)
}

func NewService() Service {
	return &service{}
}

type service struct {
}

func (s service) GetRoleBindingBySubject(subject v1Role.Subject) ([]v1Role.Banding, error) {
	db := server.DB()
	query := db.Select(customStorm.Containers("Subjects", subject))
	var rbs []v1Role.Banding
	if err := query.Find(&rbs); err != nil {
		return rbs, err
	}
	return rbs, nil
}
