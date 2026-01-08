package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type APIError struct {
	StatusCode int
	Message    string
	Internal   error
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
			var apiErr *APIError
			if errors.As(err, &apiErr) {
				ctx.JSON(apiErr.StatusCode, gin.H{"error": apiErr.Message})
				return
			}

			// Handle common DB errors
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				switch pgErr.Code {
				case "23505": // unique_violation
					ctx.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
					return
				case "23503": // foreign_key_violation
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "related record not found"})
					return
				case "23502": // not_null_violation
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing required field"})
					return
				case "23514": // check_violation
					ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data format"})
					return
				}
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}

func (e *APIError) Error() string {
	if e.Internal != nil {
		return e.Internal.Error()
	}
	return e.Message
}

func NewAPIError(statusCode int, message string, err error) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Internal:   err,
	}
}
