package postgres

import (
	"context"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type dbMigrator struct {
	migrations *migrate.Migrations
	db         *bun.DB
}

func NewDBMigrator(db *bun.DB) *dbMigrator {
	return &dbMigrator{
		migrations: migrate.NewMigrations(),
		db:         db,
	}
}

// Migrate - start db migration with list of models
func (m *dbMigrator) Migrate() error {
	log.Info("start db migration")

	// add a new model to migrate here
	models := []interface{}{
		(*repository.Report)(nil),
		(*repository.Repository)(nil),
	}

	if err := m.createTablesIfNotExist(models); err != nil {
		return err
	}

	log.Info("db migration complete")
	return nil
}

// createTablesIfNotExist - Iterate through the list of model to create new tables if doesn't exist
func (m *dbMigrator) createTablesIfNotExist(models []interface{}) error {
	for _, model := range models {
		_, err := m.db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
