package database

import (
	"github.com/marktrs/gitsast/app"
	"github.com/urfave/cli/v2"
)

func NewDBCommand() *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					_, app, err := app.StartFromCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					return NewDBMigrator(app.DB()).Migrate()
				},
			},
			{
				Name:  "migrate tables",
				Usage: "migrate tables",
				Action: func(c *cli.Context) error {
					_, app, err := app.StartFromCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					return NewDBMigrator(app.DB()).CreateTablesIfNotExist()
				},
			},
			{
				Name:  "drop-tables",
				Usage: "drop tables",
				Action: func(c *cli.Context) error {
					_, app, err := app.StartFromCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					return NewDBMigrator(app.DB()).DropTable()
				},
			},
			{
				Name:  "init-rules",
				Usage: "initialize rules",
				Action: func(c *cli.Context) error {
					_, app, err := app.StartFromCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					return NewDBMigrator(app.DB()).DropTable()
				},
			},
		},
	}
}
