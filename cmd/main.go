package main

import (
	"os"

	"github.com/marktrs/gitsast/cmd/api"
	"github.com/marktrs/gitsast/cmd/database"
	_ "github.com/marktrs/gitsast/internal/model"
	_ "github.com/marktrs/gitsast/internal/repository"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "GitSAST",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			api.NewAPICommand(),
			database.NewDBCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().AnErr("error", err).Msg("failed to run app")
	}
}
