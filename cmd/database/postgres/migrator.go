package postgres

import (
	"context"
	"database/sql"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/internal/model"
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
	log.Info("starting db migration")

	// add a new model to migrate here
	models := []interface{}{
		(*model.Rule)(nil),
		(*model.Report)(nil),
		(*model.Repository)(nil),
	}

	if err := m.resetTable(models); err != nil {
		log.Error(err)
		return err
	}

	if err := m.createTablesIfNotExist(models); err != nil {
		log.Error(err)
		return err
	}

	if err := m.insertInitialRulesIfNotExist([]*model.Rule{
		{
			ID:          1,
			Name:        "Public key leak",
			Description: "A secret starts with the prefix public_key",
			Keyword:     "public_key",
			Severity:    model.Low,
		},
		{
			ID:          2,
			Name:        "Private key leak",
			Description: "A secret starts with the prefix private_key",
			Keyword:     "private_key",
			Severity:    model.High,
		},
	}); err != nil {
		log.Error(err)
		return err
	}

	log.Info("db migration complete")
	return nil
}

// ResetTable - reset table
func (m *dbMigrator) resetTable(models []interface{}) error {
	log.Info("resetting table")

	for _, model := range models {
		_, err := m.db.NewDropTable().
			Model(model).
			IfExists().
			Exec(context.Background())
		if err != nil {
			log.Error(err)
			return err
		}
	}

	log.Info("reset table complete")
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
			log.Error(err)
			return err
		}
	}

	return nil
}

func (m *dbMigrator) insertInitialRulesIfNotExist(rules []*model.Rule) error {
	keywords := make([]string, len(rules))
	for i, rule := range rules {
		keywords[i] = rule.Keyword
	}

	return m.db.RunInTx(context.Background(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// check if exists
		exists, err := m.db.NewSelect().
			Model((*model.Rule)(nil)).
			Where("keyword IN (?)", bun.In(keywords)).
			Exists(ctx)
		if err != nil {
			log.Error(err)
			return err
		}

		if exists {
			return nil
		}

		// insert if not exists
		_, err = m.db.NewInsert().
			Model(&rules).
			Exec(context.Background())
		if err != nil {
			log.Error(err)
			return err
		}

		return nil
	})
}
