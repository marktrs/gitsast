package middleware

import (
	"net/http"

	"github.com/uptrace/bunrouter"
)

func Cors(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		origin := req.Header.Get("Origin")
		if origin == "" {
			return next(w, req)
		}

		h := w.Header()

		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Credentials", "true")

		// CORS preflight.
		if req.Method == http.MethodOptions {
			h.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD")
			h.Set("Access-Control-Allow-Headers", "authorization,content-type")
			h.Set("Access-Control-Max-Age", "86400")
			return nil
		}

		return next(w, req)
	}
}
