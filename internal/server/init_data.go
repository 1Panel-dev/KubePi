package server

import (
	"errors"
	"fmt"
	v1 "github.com/KubeOperator/ekko/internal/model/v1"
	"github.com/KubeOperator/ekko/internal/model/v1/group"
	"github.com/KubeOperator/ekko/internal/model/v1/role"
	"github.com/KubeOperator/ekko/internal/model/v1/user"
	"github.com/asdine/storm/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (e *EkkoSerer) setDefaultRoles() {
	defaultRoles := []string{"administrator", "anonymous", "authenticated", "viewer"}
	tx, err := e.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("can not start a transaction,please check your db connection: %s", err.Error()))
	}
	for i := range defaultRoles {
		n := defaultRoles[i]
		var r role.Role
		if err := tx.One("Name", n, &r); err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				_ = tx.Rollback()
				panic(fmt.Sprintf("can not query role please check db connection: %s", err.Error()))
			}
		}
		if r.Name == "" {
			r := role.Role{
				BaseModel: v1.BaseModel{
					ApiVersion: "v1",
					Kind:       "Role",
					CreateAt:   time.Now(),
					UpdateAt:   time.Now(),
					CreatedBy:  "system",
				},
				Metadata: v1.Metadata{
					Name: n,
					UUID: uuid.New().String(),
				},
			}
			if err := tx.Save(&r); err != nil {
				_ = tx.Rollback()
				panic(fmt.Sprintf("can not save role please check db connection: %s", err.Error()))
			}
			e.logger.Infof("create default role  %s ", n)
		}
	}
	_ = tx.Commit()
}

func (e *EkkoSerer) setDefaultUserGroups() {
	defaultGroupNames := []string{"administrators"}
	tx, err := e.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("can not start a transaction,please check your db connection: %s", err.Error()))
	}

	for i := range defaultGroupNames {
		n := defaultGroupNames[i]
		var g group.Group
		if err := tx.One("Name", n, &g); err != nil {
			if !errors.Is(err, storm.ErrNotFound) {
				_ = tx.Rollback()
				panic(fmt.Sprintf("can not query  user group please check db connection: %s", err.Error()))
			}
		}
		if g.Name == "" {
			g := group.Group{
				BaseModel: v1.BaseModel{
					ApiVersion: "v1",
					Kind:       "Group",
					CreateAt:   time.Now(),
					UpdateAt:   time.Now(),
					CreatedBy:  "system",
				},
				Metadata: v1.Metadata{
					Name: n,
					UUID: uuid.New().String(),
				},
			}
			if err := tx.Save(&g); err != nil {
				_ = tx.Rollback()
				panic(fmt.Sprintf("can not save user group please check db connection: %s", err.Error()))
			}
			e.logger.Infof("create default user group  %s ", n)
		}
	}
	_ = tx.Commit()

}

func (e *EkkoSerer) setSuperUser() {
	tx, err := e.db.Begin(true)
	if err != nil {
		panic(fmt.Sprintf("can not start a transaction,please check your db connection: %s", err.Error()))
	}

	var superUser user.User

	if err := tx.One("Name", "admin", &superUser); err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			_ = tx.Rollback()
			panic(fmt.Sprintf("can not query supper user please check db connection: %s", err.Error()))
		}
	}

	if superUser.Name == "" {
		e.logger.Info("creat supper user")
		superUser = user.User{
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
			Spec: user.Spec{
				Info: user.Info{
					NickName: "administrator",
					Email:    "support@fit2cloud.com",
					Language: "zh-CN",
				},
			},
		}
		pass := "admin123"
		hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost) //加密处理
		superUser.Spec.Authenticate.Password = string(hash)
		if err := tx.Save(&superUser); err != nil {
			_ = tx.Rollback()
			panic("can not save supper user please check db connection")
		}
		_ = tx.Commit()
		e.logger.Info("create supper user success")
		e.logger.Info("username: admin")
		e.logger.Info("password: admin123")
	}
}
