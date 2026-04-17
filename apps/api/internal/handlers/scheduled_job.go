package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"gorm.io/gorm"
)

type ScheduledJobHandler struct {
	db *gorm.DB
}

func NewScheduledJobHandler(db *gorm.DB) *ScheduledJobHandler {
	return &ScheduledJobHandler{db: db}
}

func (h *ScheduledJobHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query := h.db.Model(&models.ScheduledJob{})
	if keyword != "" {
		query = query.Where("name LIKE ? OR command LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.ScheduledJob
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

func (h *ScheduledJobHandler) Create(c *gin.Context) {
	var req models.CreateScheduledJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}

	entity := models.ScheduledJob{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Command:  req.Command,
		Status:   status,
		Remark:   req.Remark,
	}
	if err := h.db.Create(&entity).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("任务创建失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *ScheduledJobHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateScheduledJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.CronExpr != "" {
		updates["cron_expr"] = req.CronExpr
	}
	if req.Command != "" {
		updates["command"] = req.Command
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

	if err := h.db.Model(&models.ScheduledJob{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新失败"))
		return
	}

	var entity models.ScheduledJob
	if err := h.db.First(&entity, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("任务不存在"))
		return
	}
	c.JSON(http.StatusOK, response.Success(entity))
}

func (h *ScheduledJobHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.ScheduledJob{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}

func (h *ScheduledJobHandler) RunOnce(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var job models.ScheduledJob
	if err := h.db.First(&job, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("任务不存在"))
		return
	}

	now := time.Now()
	result := "success"
	message := "manual trigger completed"
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.ScheduledJob{}).Where("id = ?", id).Updates(map[string]interface{}{
			"last_run_at": now,
			"last_result": result,
		}).Error; err != nil {
			return err
		}

		log := models.ScheduledJobLog{
			JobID:       job.ID,
			TriggerType: "manual",
			Status:      result,
			Message:     message,
			RunAt:       now,
		}
		return tx.Create(&log).Error
	}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("执行失败"))
		return
	}

	var refreshed models.ScheduledJob
	if err := h.db.First(&refreshed, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("读取任务失败"))
		return
	}
	c.JSON(http.StatusOK, response.Success(refreshed))
}

func (h *ScheduledJobHandler) ListLogs(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	offset := (page - 1) * size
	jobID, hasJobID := parseOptionalInt(c.Query("jobId"))

	query := h.db.Model(&models.ScheduledJobLog{})
	if hasJobID {
		query = query.Where("job_id = ?", jobID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var list []models.ScheduledJobLog
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
