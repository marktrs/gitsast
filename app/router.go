package app

import (
	"github.com/labstack/gommon/log"
	"github.com/marktrs/gitsast/app/middleware"
	"github.com/uptrace/bunrouter"
)

func (app *App) initRouter() {
	log.Info("registering routes")
	r := bunrouter.New(
		bunrouter.Use(middleware.Cors),
		bunrouter.Use(middleware.WriteResponse),
		bunrouter.Use(middleware.RequestLogger),
		bunrouter.Use(middleware.ErrorHandler),
	)

	app.apiRouter = r.NewGroup("/api").NewGroup("/v1")

	log.Info("registered routes")
}
