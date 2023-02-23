package database

import (
	"context"
	"database/sql"

	"github.com/marktrs/gitsast/internal/model"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("starting db migration")
	if err := m.DropTable(); err != nil {
		log.Err(err)
		return err
	}

	if err := m.CreateTablesIfNotExist(); err != nil {
		log.Err(err)
		return err
	}

	if err := m.InsertInitialRulesIfNotExist(); err != nil {
		log.Err(err)
		return err
	}

	log.Info().Msg("db migration complete")
	return nil
}

// ResetTable - reset table
func (m *dbMigrator) DropTable() error {
	log.Info().Msg("resetting table")

	for _, model := range m.models {
		_, err := m.db.NewDropTable().
			Model(model).
			IfExists().
			Exec(context.Background())
		if err != nil {
			log.Err(err)
			return err
		}
	}

	log.Info().Msg("reset table complete")
	return nil
}

// createTablesIfNotExist - Iterate through the list of model to create new tables if doesn't exist
func (m *dbMigrator) CreateTablesIfNotExist() error {
	log.Info().Msg("creating tables if not exist")
	for _, model := range m.models {
		_, err := m.db.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(context.Background())
		if err != nil {
			log.Err(err)
			return err
		}
	}
	log.Info().Msg("created tables")
	return nil
}

func (m *dbMigrator) InsertInitialRulesIfNotExist() error {
	log.Info().Msg("initializing rules")
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
			log.Err(err)
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
			log.Err(err)
			return err
		}

		log.Info().Msg("initialized rules")

		return nil
	})
}
