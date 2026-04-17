package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type RoleHandler struct {
	db *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

func (h *RoleHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query := h.db.Model(&models.Role{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR role_key LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.Role
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

func (h *RoleHandler) Create(c *gin.Context) {
	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	dataScope := req.DataScope
	if dataScope == "" {
		dataScope = "self"
	}

	entity := models.Role{
		Name:        req.Name,
		RoleKey:     req.RoleKey,
		Permissions: req.Permissions,
		DataScope:   dataScope,
		Status:      status,
		Remark:      req.Remark,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("角色创建失败，标识可能重复"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.RoleKey != "" {
		updates["role_key"] = req.RoleKey
	}
	if req.Permissions != "" {
		updates["permissions"] = req.Permissions
	}
	if req.DataScope != "" {
		updates["data_scope"] = req.DataScope
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("no fields to update"))
		return
	}

	if err := h.db.Model(&models.Role{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("更新失败"))
		return
	}

	var entity models.Role
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("角色不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}
