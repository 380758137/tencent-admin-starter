package handlers

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/middleware"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/utils"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) List(c *gin.Context) {
	page := parseIntOr(c.Query("page"), 1)
	size := parseIntOr(c.Query("size"), 10)
	keyword := c.Query("keyword")
	offset := (page - 1) * size

	query, err := h.scopedUserQuery(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	var users []models.User
	if err := query.Order("id DESC").Offset(offset).Limit(size).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("查询失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.PageData{
		List:  users,
		Total: total,
		Page:  page,
		Size:  size,
	}))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	password := req.Password
	if password == "" {
		password = "Admin@123456"
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("密码处理失败"))
		return
	}

	role := req.Role
	role = normalizeUserRoles(role)
	if role == "" {
		role = "operator"
	}

	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	deptID := uint(0)
	if req.DeptID != nil {
		deptID = *req.DeptID
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: hash,
		DisplayName:  req.DisplayName,
		Role:         role,
		DeptID:       deptID,
		Status:       status,
	}
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, response.Error("用户创建失败，用户名可能重复"))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid payload"))
		return
	}

	updates := map[string]interface{}{}
	if req.DisplayName != "" {
		updates["display_name"] = req.DisplayName
	}
	if req.Role != "" {
		updates["role"] = normalizeUserRoles(req.Role)
	}
	if req.DeptID != nil {
		updates["dept_id"] = *req.DeptID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		hash, hashErr := utils.HashPassword(req.Password)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, response.Error("密码处理失败"))
			return
		}
		updates["password_hash"] = hash
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, response.Error("no fields to update"))
		return
	}

	if err := h.db.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("更新失败"))
		return
	}

	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, response.Error("用户不存在"))
		return
	}

	c.JSON(http.StatusOK, response.Success(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	if err := h.db.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("删除失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{"id": id}))
}

func (h *UserHandler) ExportCSV(c *gin.Context) {
	query, err := h.scopedUserQuery(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(err.Error()))
		return
	}

	var users []models.User
	if err := query.Order("id DESC").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("导出失败"))
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", "attachment; filename=users.csv")
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	if err := writer.Write([]string{"username", "displayName", "role", "status", "deptId"}); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("导出失败"))
		return
	}
	for _, u := range users {
		row := []string{u.Username, u.DisplayName, u.Role, strconv.Itoa(u.Status), strconv.FormatUint(uint64(u.DeptID), 10)}
		if err := writer.Write(row); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error("导出失败"))
			return
		}
	}
}

func (h *UserHandler) ImportCSV(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("缺少导入文件"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("读取文件失败"))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("CSV 解析失败"))
		return
	}
	if len(rows) <= 1 {
		c.JSON(http.StatusBadRequest, response.Error("CSV 内容为空"))
		return
	}

	imported := 0
	skipped := 0
	failed := 0
	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		if len(row) < 2 {
			failed++
			continue
		}
		username := strings.TrimSpace(row[0])
		displayName := strings.TrimSpace(row[1])
		if username == "" || displayName == "" {
			failed++
			continue
		}

		var count int64
		if err := h.db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
			c.JSON(http.StatusInternalServerError, response.Error("导入失败"))
			return
		}
		if count > 0 {
			skipped++
			continue
		}

		role := "operator"
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			role = strings.TrimSpace(row[2])
		}
		role = normalizeUserRoles(role)
		status := 1
		if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
			if parsed, parseErr := strconv.Atoi(strings.TrimSpace(row[3])); parseErr == nil {
				status = parsed
			}
		}
		deptID := uint(0)
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			if parsed, parseErr := strconv.ParseUint(strings.TrimSpace(row[4]), 10, 64); parseErr == nil {
				deptID = uint(parsed)
			}
		}
		password := "Admin@123456"
		if len(row) > 5 && strings.TrimSpace(row[5]) != "" {
			password = strings.TrimSpace(row[5])
		}
		hash, hashErr := utils.HashPassword(password)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, response.Error("导入失败"))
			return
		}

		user := models.User{
			Username:     username,
			DisplayName:  displayName,
			Role:         role,
			Status:       status,
			DeptID:       deptID,
			PasswordHash: hash,
		}
		if err := h.db.Create(&user).Error; err != nil {
			failed++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"imported": imported,
		"skipped":  skipped,
		"failed":   failed,
	}))
}

func (h *UserHandler) scopedUserQuery(c *gin.Context) (*gorm.DB, error) {
	query := h.db.Model(&models.User{})
	claims, ok := middleware.GetClaims(c)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	if utils.HasRole(claims.Role, "admin") {
		return query, nil
	}

	var current models.User
	if err := h.db.Select("id", "dept_id").First(&current, claims.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}

	dataScope := "self"
	roleKeys := utils.SplitRoleKeys(claims.Role)
	if len(roleKeys) > 0 {
		var roles []models.Role
		if err := h.db.Select("data_scope", "status").Where("role_key IN ?", roleKeys).Find(&roles).Error; err == nil {
			for _, role := range roles {
				if role.Status != 1 {
					continue
				}
				dataScope = utils.MergeDataScope(dataScope, role.DataScope)
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("role scope read failed")
		}
	}

	switch dataScope {
	case "all":
		return query, nil
	case "dept":
		if current.DeptID > 0 {
			deptIDs, collectErr := collectDepartmentSubtreeIDs(h.db, current.DeptID)
			if collectErr != nil {
				return nil, fmt.Errorf("dept scope read failed")
			}
			if len(deptIDs) == 0 {
				query = query.Where("id = ?", claims.UserID)
			} else {
				query = query.Where("dept_id IN ?", deptIDs)
			}
		} else {
			query = query.Where("id = ?", claims.UserID)
		}
	default:
		query = query.Where("id = ?", claims.UserID)
	}
	return query, nil
}

func normalizeUserRoles(raw string) string {
	roleKeys := utils.SplitRoleKeys(raw)
	return utils.JoinRoleKeys(roleKeys)
}
