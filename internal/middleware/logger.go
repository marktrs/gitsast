package middleware

import (
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/uptrace/bunrouter"
)

func RequestLogger(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// TODO: inject request id in logger
		log.Infof("%s %s", req.Method, req.RequestURI)

		return next(w, req)
	}
}
