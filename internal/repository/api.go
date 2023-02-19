package repository

import (
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"github.com/uptrace/bunrouter"
)

func RegisterHandlers(rg *bunrouter.Group, db *bun.DB, v *validator.Validate) {
	h := NewRepositoryHandler(db, v)
	rg.GET("/repository/:id", h.GetById)
	rg.GET("/repository", h.List)
	rg.POST("/repository", h.Add)
	rg.PUT("/repository/:id", h.Update)
	rg.DELETE("/repository/:id", h.Remove)
}
