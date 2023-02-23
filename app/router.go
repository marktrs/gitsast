package app

import (
	"net/http"

	"github.com/marktrs/gitsast/app/middleware"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
)

func (app *App) initRouter() {
	r := bunrouter.New(
		bunrouter.Use(middleware.Cors),
		bunrouter.Use(middleware.WriteResponse),
		bunrouter.Use(middleware.RequestLogger),
		bunrouter.Use(middleware.ErrorHandler),
	)

	r.GET("/health", Check)

	app.apiRouter = r.NewGroup("/api").NewGroup("/v1")
	app.router = r
	log.Info().Msg("initialized routes")
}

func Check(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"status": "ok",
	})
}
