package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type DepartmentHandler struct {
	db *gorm.DB
}

func NewDepartmentHandler(db *gorm.DB) *DepartmentHandler {
	return &DepartmentHandler{db: db}
}

func (h *DepartmentHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query := h.db.Model(&models.Department{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.Department
	if err := query.Order("parent_id ASC").Order("id ASC").Offset(offset).Limit(size).Find(&list).Error; err != nil {
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

func (h *DepartmentHandler) Create(c *gin.Context) {
	var req models.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	parentID := uint(0)
	if req.ParentID != nil {
		parentID = *req.ParentID
	}
	if err := h.validateDepartmentParent(0, parentID); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	dept := models.Department{
		Name:     req.Name,
		Code:     req.Code,
		ParentID: parentID,
		Manager:  req.Manager,
		Status:   status,
	}
	if err := h.db.Create(&dept).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("创建失败，编码可能重复"))
		return
	}

	c.JSON(http.StatusOK, response.Success(dept))
}

func (h *DepartmentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}
	var current models.Department
	if err := h.db.First(&current, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("部门不存在"))
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.ParentID != nil {
		if err := h.validateDepartmentParent(uint(id), *req.ParentID); err != nil {
			c.JSON(http.StatusBadRequest, response.Error(err.Error()))
			return
		}
		updates["parent_id"] = *req.ParentID
	}
	if req.Manager != "" {
		updates["manager"] = req.Manager
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("no fields to update"))
		return
	}

	if err := h.db.Model(&models.Department{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新失败"))
		return
	}

	var dept models.Department
	if err := h.db.First(&dept, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("部门不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(dept))
}

func (h *DepartmentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var childCount int64
	if err := h.db.Model(&models.Department{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, response.Error("存在子部门，无法删除"))
		return
	}

	var userCount int64
	if err := h.db.Model(&models.User{}).Where("dept_id = ?", id).Count(&userCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	if userCount > 0 {
		c.JSON(http.StatusBadRequest, response.Error("部门下存在用户，无法删除"))
		return
	}

	if err := h.db.Delete(&models.Department{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}

func (h *DepartmentHandler) validateDepartmentParent(currentID, parentID uint) error {
	if parentID == 0 {
		return nil
	}
	if currentID > 0 && parentID == currentID {
		return fmt.Errorf("上级部门不能是当前部门")
	}

	var parent models.Department
	if err := h.db.Select("id").First(&parent, parentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("上级部门不存在")
		}
		return fmt.Errorf("校验上级部门失败")
	}
	if currentID == 0 {
		return nil
	}

	subtreeIDs, err := collectDepartmentSubtreeIDs(h.db, currentID)
	if err != nil {
		return fmt.Errorf("校验部门层级失败")
	}
	for _, id := range subtreeIDs {
		if id == parentID {
			return fmt.Errorf("上级部门不能是当前部门或其子部门")
		}
	}

	return nil
}
