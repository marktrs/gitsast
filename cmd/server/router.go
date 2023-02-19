package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/marktrs/gitsast/internal/config"
	"github.com/marktrs/gitsast/internal/middleware"
	"github.com/marktrs/gitsast/internal/repository"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

// Routing - Add middleware to router and register HTTP handler of domains
func Routing(cfg *config.AppConfig, db *bun.DB) http.Handler {
	v := validator.New()
	r := bunrouter.New(
		bunrouter.Use(middleware.CorsMiddleware),
		bunrouter.Use(middleware.WriteResponse),
		bunrouter.Use(middleware.ErrorHandler),
		bunrouter.Use(middleware.RequestLogger),
	)

	rg := r.NewGroup("/api")
	v1 := rg.NewGroup("/v1")
	repository.RegisterHandlers(v1, db, v)

	return r
}
