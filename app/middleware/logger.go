package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bunrouter"
)

type loggerCtxKey struct{}

func RequestLogger(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		requestID := uuid.New().String()

		ctx := req.Context()
		ctx = context.WithValue(ctx, loggerCtxKey{}, requestID)
		req = req.WithContext(ctx)

		log.Info().Fields(map[string]interface{}{
			"request_id": requestID,
			"method":     req.Method,
			"path":       req.RequestURI,
		}).Msg("request received")

		w.Header().Set("X-Request-Id", requestID)

		return next(w, req)
	}
}
