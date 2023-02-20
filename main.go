package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/cmd/database/postgres"
	"github.com/marktrs/gitsast/cmd/server"
	"github.com/marktrs/gitsast/internal/config"
	"github.com/marktrs/gitsast/internal/queue"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var configPath = flag.String("config", "./config/local.yaml", "path to the config file")

func main() {
	// load app config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// connect db
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DB.Dsn))), pgdialect.New())

	// TODO: only start migration if flag is true, move to db CLI command
	// var isStartMigration = flag.Bool("migrate", false, "start migration")
	// if *isStartMigration {}
	// migrate db
	migrator := postgres.NewDBMigrator(db)
	if err := migrator.Migrate(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// setup queue handler
	q := queue.NewHandler()
	if err := q.StartConsumer(); err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// setup routing
	r := server.Routing(cfg, db, q)

	// start api server
	if err := server.Start(cfg, r); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
