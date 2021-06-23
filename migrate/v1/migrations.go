package v1

import (
	"fmt"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	v1Role "github.com/KubeOperator/ekko/internal/model/v1/role"
	v1User "github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/KubeOperator/ekko/migrate/migrations"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var Migrations = []migrations.Migration{
	CreateAdministrator,
}

// 创建默认角色: administrator/authenticated/anonymous/viewer
// 创建默认用户：admin
// 关联: admin 和 administrator
var CreateAdministrator = migrations.Migration{
	Version: 1,
	Message: "CreateAdministrator",
	Handler: func(db storm.Node) error {
		roleAdministrator := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				CreatedBy:  "system",
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:     "administrator",
				Description: "_administrator",
				UUID:     uuid.New().String(),
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"*"},
					Verbs:    []string{"*"},
				},
			},
		}
		defaultPass := "admin123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)
		userAdmin := v1User.User{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "User",
				CreatedBy:  "system",
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name: "admin",
				UUID: uuid.New().String(),
			},
			Spec: v1User.Spec{
				Info: v1User.Info{
					NickName: "administrator",
					Email:    "support@fit2cloud.com",
					Language: "zh-CN",
				},
				Authenticate: v1User.Authenticate{
					Password: string(hash),
					Token:    "",
				},
			},
		}

		binding := v1Role.Binding{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "RoleBinding",
				CreatedBy:  "system",
				CreateAt:   time.Now(),
				UpdateAt:   time.Now()},
			Metadata: v1.Metadata{
				Name: fmt.Sprintf("role-binding-%s-%s", roleAdministrator.Name, userAdmin.Name),
				UUID: uuid.New().String(),
			},
			Subjects: []v1Role.Subject{
				{
					Kind: "User",
					Name: userAdmin.Name,
				},
			},
			RoleRef: roleAdministrator.Name,
		}
		if err := db.Save(&roleAdministrator); err != nil {
			return err
		}
		if err := db.Save(&userAdmin); err != nil {
			return err
		}
		if err := db.Save(&binding); err != nil {
			return err
		}
		return nil
	},
}
