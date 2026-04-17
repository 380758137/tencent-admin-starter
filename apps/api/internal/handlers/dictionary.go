package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type DictionaryHandler struct {
	db *gorm.DB
}

func NewDictionaryHandler(db *gorm.DB) *DictionaryHandler {
	return &DictionaryHandler{db: db}
}

func (h *DictionaryHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	dictType := c.Query("dictType")
	offset := (page - 1) * size

	query := h.db.Model(&models.DictionaryItem{})
	if keyword != "" {
		query = query.Where("dict_type LIKE ? OR dict_label LIKE ? OR dict_value LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.DictionaryItem
	if err := query.Order("dict_type ASC").Order("sort ASC").Order("id ASC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
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

func (h *DictionaryHandler) Create(c *gin.Context) {
	var req models.CreateDictionaryItemRequest
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

	entity := models.DictionaryItem{
		DictType:  req.DictType,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		Sort:      sort,
		Status:    status,
		Remark:    req.Remark,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("字典项创建失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *DictionaryHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateDictionaryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.DictType != "" {
		updates["dict_type"] = req.DictType
	}
	if req.DictLabel != "" {
		updates["dict_label"] = req.DictLabel
	}
	if req.DictValue != "" {
		updates["dict_value"] = req.DictValue
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

	if err := h.db.Model(&models.DictionaryItem{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新失败"))
		return
	}

	var entity models.DictionaryItem
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("字典项不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *DictionaryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.DictionaryItem{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}
