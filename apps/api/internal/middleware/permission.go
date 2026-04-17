package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

func RequirePerms(db *gorm.DB, perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("unauthorized"))
			return
		}
		if claims.Role == "admin" {
			c.Next()
			return
		}

		var role models.Role
		if err := db.Select("permissions", "status").Where("role_key = ?", claims.Role).First(&role).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
			return
		}
		if role.Status != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
			return
		}

		allow := map[string]struct{}{}
		for _, p := range splitPerms(role.Permissions) {
			allow[p] = struct{}{}
		}
		if _, exists := allow["*"]; exists {
			c.Next()
			return
		}
		for _, p := range perms {
			if _, exists := allow[p]; exists {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
	}
}

func splitPerms(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
