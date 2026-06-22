package server

import (
	"github.com/1Panel-dev/webkubectl/gotty/webtty"
)

// Slave is webtty.Slave with some additional methods.
type Slave interface {
	webtty.Slave

	Close() error
}

type Factory interface {
	Name() string
	New(params map[string][]string) (Slave, error)
}
