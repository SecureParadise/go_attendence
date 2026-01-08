package middleware

import (
	"fmt"
	"net/http"

	"github.com/SecureParadise/go_attendence/internal/auth"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(AuthorizationPayloadKey)
		if !exists {
			ctx.Error(NewAPIError(http.StatusUnauthorized, "authorization payload not found", nil))
			ctx.Abort()
			return
		}

		authPayload, ok := payload.(*auth.Payload)
		if !ok {
			ctx.Error(NewAPIError(http.StatusInternalServerError, "invalid authorization payload", nil))
			ctx.Abort()
			return
		}

		hasRole := false
		for _, role := range roles {
			if authPayload.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			ctx.Error(NewAPIError(http.StatusForbidden, fmt.Sprintf("permission denied for role: %s", authPayload.Role), nil))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func HierarchicalRoleMiddleware(requiredRole string) gin.HandlerFunc {
	roleHierarchy := map[string]int{
		"admin":   4, // Global Admin access
		"hod":     3,
		"teacher": 2,
		"student": 1,
	}

	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(AuthorizationPayloadKey)
		if !exists {
			ctx.Error(NewAPIError(http.StatusUnauthorized, "authorization payload not found", nil))
			ctx.Abort()
			return
		}

		authPayload, ok := payload.(*auth.Payload)
		if !ok {
			ctx.Error(NewAPIError(http.StatusInternalServerError, "invalid authorization payload", nil))
			ctx.Abort()
			return
		}

		userRole := authPayload.Role
		if roleHierarchy[userRole] < roleHierarchy[requiredRole] {
			ctx.Error(NewAPIError(http.StatusForbidden, fmt.Sprintf("insufficient permissions: required %s, got %s", requiredRole, userRole), nil))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
