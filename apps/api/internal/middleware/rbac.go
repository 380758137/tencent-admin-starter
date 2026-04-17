package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	allow := map[string]struct{}{}
	for _, r := range roles {
		allow[r] = struct{}{}
	}

	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("unauthorized"))
			return
		}

		if _, exists := allow[claims.Role]; !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
			return
		}

		c.Next()
	}
}
