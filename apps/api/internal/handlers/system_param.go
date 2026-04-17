package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type SystemParamHandler struct {
	db *gorm.DB
}

func NewSystemParamHandler(db *gorm.DB) *SystemParamHandler {
	return &SystemParamHandler{db: db}
}

func (h *SystemParamHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query := h.db.Model(&models.SystemParam{})
	if keyword != "" {
		query = query.Where("param_key LIKE ? OR param_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.SystemParam
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

func (h *SystemParamHandler) Create(c *gin.Context) {
	var req models.CreateSystemParamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	entity := models.SystemParam{
		ParamKey:   req.ParamKey,
		ParamValue: req.ParamValue,
		ParamName:  req.ParamName,
		Status:     status,
		Remark:     req.Remark,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("参数创建失败，键值可能重复"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *SystemParamHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateSystemParamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.ParamValue != "" {
		updates["param_value"] = req.ParamValue
	}
	if req.ParamName != "" {
		updates["param_name"] = req.ParamName
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

	if err := h.db.Model(&models.SystemParam{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新失败"))
		return
	}

	var entity models.SystemParam
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("参数不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *SystemParamHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.SystemParam{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}
