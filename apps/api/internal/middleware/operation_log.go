package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"gorm.io/gorm"
)

func OperationLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		if !strings.HasPrefix(path, "/api/") {
			return
		}

		var userID uint
		var username string
		if claims, ok := GetClaims(c); ok {
			userID = claims.UserID
			username = claims.Username
		}

		entry := models.OperationLog{
			UserID:     userID,
			Username:   username,
			Method:     c.Request.Method,
			Path:       path,
			StatusCode: c.Writer.Status(),
			LatencyMs:  time.Since(start).Milliseconds(),
			IP:         c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}
		if err := db.Create(&entry).Error; err != nil {
			_ = c.Error(err)
		}
	}
}
