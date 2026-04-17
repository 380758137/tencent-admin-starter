package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type OnlineUserHandler struct {
	db *gorm.DB
}

type onlineUserRow struct {
	Username     string `json:"username"`
	LastActiveAt string `json:"lastActiveAt"`
	RequestCount int64  `json:"requestCount"`
}

func NewOnlineUserHandler(db *gorm.DB) *OnlineUserHandler {
	return &OnlineUserHandler{db: db}
}

func (h *OnlineUserHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size
	activeSince := time.Now().Add(-30 * time.Minute)

	totalQuery := h.db.Table("operation_logs").Where("username <> '' AND created_at >= ?", activeSince)
	if keyword != "" {
		totalQuery = totalQuery.Where("username LIKE ?", "%"+keyword+"%")
	}
	var total int64
	if err := totalQuery.Distinct("username").Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	listQuery := h.db.
		Table("operation_logs").
		Select("username, MAX(created_at) AS last_active_at, COUNT(*) AS request_count").
		Where("username <> '' AND created_at >= ?", activeSince)
	if keyword != "" {
		listQuery = listQuery.Where("username LIKE ?", "%"+keyword+"%")
	}

	var list []onlineUserRow
	if err := listQuery.
		Group("username").
		Order("last_active_at DESC").
		Offset(offset).
		Limit(size).
		Scan(&list).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.PageData{
		List:  list,
		Total: total,
		Page:  page,
		Size:  size,
	}))
}
