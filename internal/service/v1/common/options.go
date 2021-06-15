package common

import (
	"github.com/KubeOperator/ekko/internal/server"
	"github.com/asdine/storm/v3"
)

type DBService interface {
	GetDB(options DBOptions) storm.Node
}

type DefaultDBService struct {
}

func (d *DefaultDBService) GetDB(options DBOptions) storm.Node {
	if options.DB != nil {
		return options.DB
	}
	return server.DB()
}

type DBOptions struct {
	DB storm.Node
}
