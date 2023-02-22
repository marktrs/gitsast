package app

import (
	"github.com/marktrs/gitsast/app/middleware"
	"github.com/uptrace/bunrouter"
)

func (app *App) initRouter() {
	r := bunrouter.New(
		bunrouter.Use(middleware.Cors),
		bunrouter.Use(middleware.WriteResponse),
		bunrouter.Use(middleware.RequestLogger),
		bunrouter.Use(middleware.ErrorHandler),
	)

	app.apiRouter = r.NewGroup("/api").NewGroup("/v1")
	app.router = r
}
