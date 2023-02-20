package repository

import (
	"github.com/go-playground/validator/v10"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/marktrs/gitsast/internal/queue"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

func RegisterHandlers(rg *bunrouter.Group, db *bun.DB, v *validator.Validate, q queue.Handler) {
	rs := model.NewRepositoryRepo(db)
	rp := model.NewReportRepo(db)
	s := NewService(v, q, rs, rp)
	h := NewHTTPHandler(s)

	rg.WithGroup("/repository", func(g *bunrouter.Group) {
		g.GET("/:id", h.GetById)
		g.GET("", h.List)
		g.POST("", h.Add)
		g.PUT("/:id", h.Update)
		g.DELETE("/:id", h.Remove)
		g.POST("/:id/scan", h.Scan)
	})
}
