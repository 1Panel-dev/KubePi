package v1

import (
	v1 "github.com/KubeOperator/kubepi/internal/model/v1"
	v1Role "github.com/KubeOperator/kubepi/internal/model/v1/role"
	v1User "github.com/KubeOperator/kubepi/internal/model/v1/user"
	"github.com/KubeOperator/kubepi/migrate/migrations"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var Migrations = []migrations.Migration{
	CreateAdministrator,
}

// 创建默认系统角色: Admin |Manage Cluster| Manage User|Read only|Common User | Manage Chart
// 创建用户
// 创建用户和角色的关联

var CreateAdministrator = migrations.Migration{
	Version: 1,
	Message: "Create default user and cluster",
	Handler: func(db storm.Node) error {
		//



		roleManageClusters := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:        "Manage Clusters",
				UUID:        uuid.New().String(),
				Description: "i18n_user_manage_cluster",
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"clusters"},
					Verbs:    []string{"*"},
				},
			},
		}
		roleManageRBAC := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:        "Manage RBAC",
				Description: "i18n_user_manage_rbac",
				UUID:        uuid.New().String(),
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"users", "roles"},
					Verbs:    []string{"*"},
				},
			},
		}
		roleReadOnly := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:        "ReadOnly",
				Description: "i18n_user_manage_readonly",
				UUID:        uuid.New().String(),
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"*"},
					Verbs:    []string{"get", "list"},
				},
			},
		}
		roleCommonUser := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:        "Common User",
				Description: "i18n_user_common_user",
				UUID:        uuid.New().String(),
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"clusters"},
					Verbs:    []string{"get", "list"},
				},
			},
		}
		roleManageChart := v1Role.Role{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "Role",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name:        "Manage Chart",
				Description: "i18n_user_manage_chart",
				UUID:        uuid.New().String(),
			},
			Rules: []v1Role.PolicyRule{
				{
					Resource: []string{"charts"},
					Verbs:    []string{"*"},
				},
			},
		}

		// 创建管理员用户
		defaultUserPass := "kubepi"
		hash, _ := bcrypt.GenerateFromPassword([]byte(defaultUserPass), bcrypt.DefaultCost)
		userAdmin := v1User.User{
			BaseModel: v1.BaseModel{
				ApiVersion: "v1",
				Kind:       "User",
				BuiltIn:    true,
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
			},
			Metadata: v1.Metadata{
				Name: "admin",
				UUID: uuid.New().String(),
			},
			IsAdmin:  true,
			NickName: "Administrator",
			Email:    "support@fit2cloud.com",
			Language: "zh-CN",
			Authenticate: v1User.Authenticate{
				Password: string(hash),
			},
		}
		// 创建绑定关系
		dbObjects := []interface{}{
			&roleManageClusters, &roleCommonUser, &roleManageRBAC, &roleReadOnly,
			&userAdmin, &roleManageChart,
		}
		for i := range dbObjects {
			if err := db.Save(dbObjects[i]); err != nil {
				return err
			}
		}
		return nil
	},
}
