package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type MonitorHandler struct {
	db *gorm.DB
}

func NewMonitorHandler(db *gorm.DB) *MonitorHandler {
	return &MonitorHandler{db: db}
}

func (h *MonitorHandler) Overview(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("数据库状态读取失败"))
		return
	}

	dbHealth := "ok"
	if pingErr := sqlDB.Ping(); pingErr != nil {
		dbHealth = pingErr.Error()
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	counts := map[string]int64{}
	modelRefs := map[string]interface{}{
		"users":         &models.User{},
		"departments":   &models.Department{},
		"roles":         &models.Role{},
		"menus":         &models.Menu{},
		"dictionary":    &models.DictionaryItem{},
		"systemParams":  &models.SystemParam{},
		"operationLogs": &models.OperationLog{},
		"loginLogs":     &models.LoginLog{},
	}
	for key, model := range modelRefs {
		var count int64
		if countErr := h.db.Model(model).Count(&count).Error; countErr != nil {
			c.JSON(http.StatusInternalServerError, response.Error("统计失败"))
			return
		}
		counts[key] = count
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"runtime": gin.H{
			"goroutines": runtime.NumGoroutine(),
			"heapAlloc":  mem.HeapAlloc,
			"heapSys":    mem.HeapSys,
			"goVersion":  runtime.Version(),
			"now":        time.Now().Format(time.RFC3339),
		},
		"database": gin.H{
			"health": dbHealth,
			"stats":  sqlDB.Stats(),
		},
		"counts": counts,
	}))
}
