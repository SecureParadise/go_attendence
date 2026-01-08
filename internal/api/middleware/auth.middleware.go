package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorization
func AuthMiddleware(tokenMaker auth.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.Error(NewAPIError(http.StatusUnauthorized, err.Error(), err))
			ctx.Abort()
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.Error(NewAPIError(http.StatusUnauthorized, err.Error(), err))
			ctx.Abort()
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.Error(NewAPIError(http.StatusUnauthorized, err.Error(), err))
			ctx.Abort()
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken, auth.AccessToken)
		if err != nil {
			ctx.Error(NewAPIError(http.StatusUnauthorized, err.Error(), err))
			ctx.Abort()
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
