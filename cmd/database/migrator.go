package database

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
	models     []interface{}
}

func NewDBMigrator(db *bun.DB) *dbMigrator {
	models := []interface{}{
		(*model.Rule)(nil),
		(*model.Report)(nil),
		(*model.Repository)(nil),
	}

	return &dbMigrator{
		migrations: migrate.NewMigrations(),
		db:         db,
		models:     models,
	}
}

// Migrate - start db migration with list of models
func (m *dbMigrator) Migrate() error {
	log.Info("starting db migration")
	if err := m.DropTable(); err != nil {
		log.Error(err)
		return err
	}

	if err := m.CreateTablesIfNotExist(); err != nil {
		log.Error(err)
		return err
	}

	if err := m.InsertInitialRulesIfNotExist(); err != nil {
		log.Error(err)
		return err
	}

	log.Info("db migration complete")
	return nil
}

// ResetTable - reset table
func (m *dbMigrator) DropTable() error {
	log.Info("resetting table")

	for _, model := range m.models {
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
func (m *dbMigrator) CreateTablesIfNotExist() error {
	for _, model := range m.models {
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

func (m *dbMigrator) InsertInitialRulesIfNotExist() error {
	rules := []*model.Rule{
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
	}

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
