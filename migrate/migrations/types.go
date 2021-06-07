package migrations

import "github.com/asdine/storm/v3"

type MigrationFUNC func(db storm.Node) error

type Migration struct {
	Version int
	Handler MigrationFUNC
	Message string
}

