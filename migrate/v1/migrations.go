package v1

import (
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/migrate/migrations"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
)

var Migrations = []migrations.Migration{
	InitData,
	CreateRole,
}

var InitData = migrations.Migration{
	Version: 1,
	Message: "init data",
	Handler: func(db storm.Node) error {
		role := v1Role.Role{
			BaseModel: v1.BaseModel{},
			Metadata: v1.Metadata{
				Name: "test1",
				UUID: uuid.New().String(),
			},
			Rules: nil,
		}
		return db.Save(&role)
	},
}

var CreateRole = migrations.Migration{
	Version: 2,
	Message: "init data 123",
	Handler: func(db storm.Node) error {
		role := v1Role.Role{
			BaseModel: v1.BaseModel{},
			Metadata: v1.Metadata{
				Name: "test2",
				UUID: uuid.New().String(),
			},
			Rules: nil,
		}
		return db.Save(&role)
	},
}
