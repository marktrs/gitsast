package repository

import (
	"context"

	"github.com/marktrs/gitsast/app"
	"github.com/marktrs/gitsast/internal/model"
	"github.com/uptrace/bunrouter"
)

func init() {
	app.OnStart("repository.initRoutes", func(ctx context.Context, app *app.App) error {
		rs := model.NewRepositoryRepo(app.DB())
		rp := model.NewReportRepo(app.DB())
		s := NewService(app, rs, rp)
		h := NewHTTPHandler(s)

		app.APIRouter().WithGroup("/repository", func(g *bunrouter.Group) {
			g.GET("/:id", h.GetById)
			g.GET("", h.List)
			g.POST("", h.Add)
			g.PUT("/:id", h.Update)
			g.DELETE("/:id", h.Remove)
			g.POST("/:id/scan", h.Scan)
			g.GET("/:id/report", h.GetReport)
		})

		return nil
	})
}
