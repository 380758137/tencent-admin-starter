package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/utils"
	"gorm.io/gorm"
)

func RequirePerms(db *gorm.DB, perms ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := GetClaims(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.Error("unauthorized"))
			return
		}
		if utils.HasRole(claims.Role, "admin") {
			c.Next()
			return
		}

		roleKeys := utils.SplitRoleKeys(claims.Role)
		if len(roleKeys) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
			return
		}

		var roles []models.Role
		if err := db.Select("permissions", "status").Where("role_key IN ?", roleKeys).Find(&roles).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, response.Error("forbidden"))
			return
		}

		allow := map[string]struct{}{}
		for _, role := range roles {
			if role.Status != 1 {
				continue
			}
			for _, p := range splitPerms(role.Permissions) {
				allow[p] = struct{}{}
			}
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
