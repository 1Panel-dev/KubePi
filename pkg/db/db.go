package db

import (
	"github.com/KubeOperator/ekko/pkg/config"
	"github.com/asdine/storm/v3"
)

var instance *storm.DB

type InitDbPhase struct {
}

func (i *InitDbPhase) Init() error {
	c := config.ReadConfig()
	di, err := storm.Open(c.Spec.Db.Path)
	if err != nil {
		return err
	}
	instance = di
	return nil
}

func (i *InitDbPhase) Name() string {
	return "db"
}

func GetDB() *storm.DB {
	return instance
}
