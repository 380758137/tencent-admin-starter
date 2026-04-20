package handlers

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/config"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/middleware"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/response"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	cfg config.Config
	db  *gorm.DB
}

func NewAuthHandler(cfg config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if logErr := h.writeLoginLog(c, req.Username, 0, "invalid login payload"); logErr != nil {
			_ = c.Error(logErr)
		}
		c.JSON(http.StatusBadRequest, response.Error("invalid login payload"))
		return
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if logErr := h.writeLoginLog(c, req.Username, 0, "user not found"); logErr != nil {
			_ = c.Error(logErr)
		}
		c.JSON(http.StatusUnauthorized, response.Error("用户名或密码错误"))
		return
	}

	if user.Status != 1 {
		if logErr := h.writeLoginLog(c, req.Username, 0, "user disabled"); logErr != nil {
			_ = c.Error(logErr)
		}
		c.JSON(http.StatusForbidden, response.Error("用户已禁用"))
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		if logErr := h.writeLoginLog(c, req.Username, 0, "wrong password"); logErr != nil {
			_ = c.Error(logErr)
		}
		c.JSON(http.StatusUnauthorized, response.Error("用户名或密码错误"))
		return
	}

	token, err := middleware.NewToken(h.cfg, user.ID, user.Username, user.Role)
	if err != nil {
		if logErr := h.writeLoginLog(c, req.Username, 0, "token issue failed"); logErr != nil {
			_ = c.Error(logErr)
		}
		c.JSON(http.StatusInternalServerError, response.Error("签发 token 失败"))
		return
	}
	if logErr := h.writeLoginLog(c, req.Username, 1, "login success"); logErr != nil {
		_ = c.Error(logErr)
	}
	perms, dataScope, permErr := h.loadRoleAccess(user.Role)
	if permErr != nil {
		c.JSON(http.StatusInternalServerError, response.Error("权限读取失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"token": token,
		"user": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"displayName": user.DisplayName,
			"role":        user.Role,
			"deptId":      user.DeptID,
			"permissions": perms,
			"dataScope":   dataScope,
		},
	}))
}

func (h *AuthHandler) writeLoginLog(c *gin.Context, username string, status int, message string) error {
	entry := models.LoginLog{
		Username:  username,
		Status:    status,
		Message:   message,
		IP:        c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}
	return h.db.Create(&entry).Error
}

func (h *AuthHandler) Me(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
		return
	}

	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("用户不存在"))
		return
	}
	perms, dataScope, permErr := h.loadRoleAccess(user.Role)
	if permErr != nil {
		c.JSON(http.StatusInternalServerError, response.Error("权限读取失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"displayName": user.DisplayName,
		"role":        user.Role,
		"deptId":      user.DeptID,
		"permissions": perms,
		"dataScope":   dataScope,
		"status":      user.Status,
	}))
}

func (h *AuthHandler) Menus(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("unauthorized"))
		return
	}

	var user models.User
	if err := h.db.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("用户不存在"))
		return
	}

	perms, _, permErr := h.loadRoleAccess(user.Role)
	if permErr != nil {
		c.JSON(http.StatusInternalServerError, response.Error("权限读取失败"))
		return
	}

	allowAll := false
	permSet := make(map[string]struct{}, len(perms))
	for _, perm := range perms {
		trimmed := strings.TrimSpace(perm)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			allowAll = true
			break
		}
		permSet[trimmed] = struct{}{}
	}

	var menus []models.Menu
	if err := h.db.
		Where("status = ? AND menu_type = ?", 1, "menu").
		Order("parent_id ASC").
		Order("sort ASC").
		Order("id ASC").
		Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("菜单读取失败"))
		return
	}

	if allowAll {
		c.JSON(http.StatusOK, response.Success(menus))
		return
	}

	filtered := make([]models.Menu, 0, len(menus))
	for _, menu := range menus {
		rawPerm := strings.TrimSpace(menu.Perms)
		if rawPerm == "" {
			filtered = append(filtered, menu)
			continue
		}
		if _, exists := permSet[rawPerm]; exists {
			filtered = append(filtered, menu)
		}
	}

	c.JSON(http.StatusOK, response.Success(filtered))
}

func (h *AuthHandler) loadRoleAccess(roleRaw string) ([]string, string, error) {
	roleKeys := utils.SplitRoleKeys(roleRaw)
	if len(roleKeys) == 0 {
		return []string{}, "self", nil
	}

	var roles []models.Role
	if err := h.db.Select("permissions", "data_scope", "status").Where("role_key IN ?", roleKeys).Find(&roles).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, "self", nil
		}
		return nil, "", err
	}

	scope := "self"
	allow := map[string]struct{}{}
	for _, role := range roles {
		if role.Status != 1 {
			continue
		}
		scope = utils.MergeDataScope(scope, role.DataScope)
		for _, part := range strings.Split(role.Permissions, ",") {
			perm := strings.TrimSpace(part)
			if perm == "" {
				continue
			}
			allow[perm] = struct{}{}
		}
	}

	if _, exists := allow["*"]; exists {
		return []string{"*"}, scope, nil
	}

	out := make([]string, 0, len(allow))
	for perm := range allow {
		out = append(out, perm)
	}
	sort.Strings(out)
	return out, scope, nil
}
