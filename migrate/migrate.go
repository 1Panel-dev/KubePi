package migrate

import (
	"errors"
	"github.com/KubeOperator/kubepi/migrate/migrations"
	v1 "github.com/KubeOperator/kubepi/migrate/v1"
	"github.com/asdine/storm/v3"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
)

var definedMigrations = append([]migrations.Migration{}, v1.Migrations...)

func RunMigrate(db *storm.DB, logger *logrus.Logger) {
	var currentDbVersion int
	if err := db.Get("db", "current_db_version", &currentDbVersion); err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			currentDbVersion = 0
		} else {
			logger.Errorf("can not get current db version ,%s", err.Error())
			os.Exit(1)
		}
	}
	logger.Infof("current db version: %d",currentDbVersion)

	var prepareExecuteMigrationVersions []int

	for i := range definedMigrations {
		if definedMigrations[i].Version > currentDbVersion {
			prepareExecuteMigrationVersions = append(prepareExecuteMigrationVersions, definedMigrations[i].Version)
		}
	}

	if len(prepareExecuteMigrationVersions) > 0 {
		sort.Ints(prepareExecuteMigrationVersions)
		lastVersion := prepareExecuteMigrationVersions[len(prepareExecuteMigrationVersions)-1]

		tx, err := db.Begin(true)
		if err != nil {
			logger.Errorf("can not open transaction ,%s", err.Error())
			os.Exit(1)
		}
		for i := range definedMigrations {
			for j := range prepareExecuteMigrationVersions {
				if definedMigrations[i].Version == prepareExecuteMigrationVersions[j] {
					logger.Infof("executing db migration: [%d]  %s", definedMigrations[i].Version, definedMigrations[i].Message)
					if err := definedMigrations[i].Handler(tx); err != nil {
						_ = tx.Rollback()
						logger.Errorf("execute migration: [%d] %s failed,rollback it", definedMigrations[i].Version, err.Error())
						os.Exit(1)
					}
				}
			}
		}
		currentDbVersion = lastVersion
		if err := tx.Set("db", "current_db_version", currentDbVersion); err != nil {
			_ = tx.Rollback()
			logger.Errorf("update db version failed %s ,rollback it", err.Error())
			os.Exit(1)
		}
		logger.Infof("update db to version: %d", currentDbVersion)
		_ = tx.Commit()
	}
}
