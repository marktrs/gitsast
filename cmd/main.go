package main

import (
	"net/http"
	"os"

	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/app"
	_ "github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/recover"
	_ "github.com/marktrs/gitsast/internal/repository"
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
			apiCommand,
			// newDBCommand(migrations.Migrations),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var apiCommand = &cli.Command{
	Name:  "api",
	Usage: "start GitSAST API server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: ":8000",
			Usage: "serve address",
		},
		&cli.StringFlag{
			Name:  "config",
			Value: "./config/dev.yaml",
			Usage: "path to environment config file",
		},
	},
	Action: func(c *cli.Context) error {
		ctx, api, err := app.StartFromCLI(c)
		if err != nil {
			return err
		}
		defer api.Stop()

		var handler http.Handler
		handler = api.Router()
		handler = recover.PanicHandler{Next: handler}

		cfg := api.Config()
		srv := &http.Server{
			Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
			Handler:      handler,
		}

		if err := api.Queue().StartConsumer(); err != nil {
			return err
		}

		go func() {
			if err := srv.ListenAndServe(); err != nil && !isServerClosed(err) {
				log.Printf("ListenAndServe failed: %s", err)
			}
		}()

		log.Printf("listening on %s", srv.Addr)
		app.WaitExitSignal()

		return srv.Shutdown(ctx)
	},
}

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}
