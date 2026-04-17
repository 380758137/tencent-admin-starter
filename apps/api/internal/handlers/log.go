package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type LogHandler struct {
	db *gorm.DB
}

func NewLogHandler(db *gorm.DB) *LogHandler {
	return &LogHandler{db: db}
}

func (h *LogHandler) ListOperations(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	username := c.Query("username")
	method := c.Query("method")
	pathKeyword := c.Query("path")
	offset := (page - 1) * size

	query := h.db.Model(&models.OperationLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if method != "" {
		query = query.Where("method = ?", method)
	}
	if pathKeyword != "" {
		query = query.Where("path LIKE ?", "%"+pathKeyword+"%")
	}
	if status, ok := parseOptionalInt(c.Query("status")); ok {
		query = query.Where("status_code = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.OperationLog
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
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

func (h *LogHandler) ListLogins(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	username := c.Query("username")
	offset := (page - 1) * size

	query := h.db.Model(&models.LoginLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status, ok := parseOptionalInt(c.Query("status")); ok {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.LoginLog
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
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
