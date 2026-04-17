package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type PositionHandler struct {
	db *gorm.DB
}

func NewPositionHandler(db *gorm.DB) *PositionHandler {
	return &PositionHandler{db: db}
}

func (h *PositionHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query := h.db.Model(&models.Position{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.Position
	if err := query.Order("sort ASC").Order("id DESC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
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

func (h *PositionHandler) Create(c *gin.Context) {
	var req models.CreatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	sort := 0
	if req.Sort != nil {
		sort = *req.Sort
	}
	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	entity := models.Position{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   sort,
		Status: status,
		Remark: req.Remark,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("岗位创建失败，编码可能重复"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *PositionHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdatePositionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
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

	if err := h.db.Model(&models.Position{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("更新失败"))
		return
	}

	var entity models.Position
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("岗位不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *PositionHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.Position{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}
