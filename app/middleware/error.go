package middleware

import (
	"database/sql"
	"io"
	"net/http"

	"github.com/labstack/gommon/log"
	"github.com/uptrace/bunrouter"
)

func ErrorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)

		if err != nil {
			log.Error(err)

			switch err := err.(type) {
			case HTTPError:
				w.WriteHeader(err.statusCode)
				_ = bunrouter.JSON(w, err)
			default:
				httpErr := NewHTTPError(err)
				w.WriteHeader(httpErr.statusCode)
				_ = bunrouter.JSON(w, httpErr)
			}

			return err
		}

		return nil
	}
}

type HTTPError struct {
	statusCode int

	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewHTTPError(err error) HTTPError {
	switch err {
	case io.EOF:
		return HTTPError{
			statusCode: http.StatusBadRequest,
			Code:       "eof",
			Message:    "EOF reading HTTP request body",
		}
	case sql.ErrNoRows:
		return HTTPError{
			statusCode: http.StatusNotFound,
			Code:       "not_found",
			Message:    "Row not found",
		}
	}

	return HTTPError{
		statusCode: http.StatusInternalServerError,
		Code:       "internal",
		Message:    "Internal server error",
	}
}
