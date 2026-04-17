package handlers

import (
	"errors"
	"net/http"
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

func (h *AuthHandler) loadRoleAccess(roleKey string) ([]string, string, error) {
	var role models.Role
	if err := h.db.Select("permissions", "data_scope").Where("role_key = ?", roleKey).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []string{}, "self", nil
		}
		return nil, "", err
	}
	raw := strings.TrimSpace(role.Permissions)
	if raw == "" {
		scope := role.DataScope
		if scope == "" {
			scope = "self"
		}
		return []string{}, scope, nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	scope := role.DataScope
	if scope == "" {
		scope = "self"
	}
	return out, scope, nil
}
