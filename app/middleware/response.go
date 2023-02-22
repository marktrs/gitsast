package middleware

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func WriteResponse(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return next(w, req)
	}
}
