package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
)

func RequestLogger(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		req.WithContext(req.Context())

		request_id := uuid.New().String()
		log.Info().Fields(map[string]interface{}{
			"method":     req.Method,
			"path":       req.RequestURI,
			"request_id": request_id,
			"params":     req.Params().Map(),
		}).Msg("request received")

		w.Header().Set("X-Request-Id", request_id)

		return next(w, req)
	}
}
