package database

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/liusheng/tencent-admin-starter/apps/api/internal/config"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/models"
	"github.com/liusheng/tencent-admin-starter/apps/api/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(cfg config.Config) (*gorm.DB, error) {
	var driver gorm.Dialector = mysql.Open(cfg.MySQLDSN)
	useSQLite := strings.HasPrefix(cfg.MySQLDSN, "sqlite://")
	if useSQLite {
		driver = sqlite.Open(strings.TrimPrefix(cfg.MySQLDSN, "sqlite://"))
	}
	db, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if !useSQLite {
		if err := runSQLMigrations(db, cfg.MigrationDir); err != nil {
			return nil, err
		}
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Role{},
		&models.Menu{},
		&models.DictionaryItem{},
		&models.SystemParam{},
		&models.OperationLog{},
		&models.LoginLog{},
		&models.Position{},
		&models.Notice{},
		&models.ScheduledJob{},
		&models.ScheduledJobLog{},
	); err != nil {
		return nil, err
	}

	if err := seedDefaultAdmin(db); err != nil {
		return nil, err
	}
	if err := seedFoundationData(db); err != nil {
		return nil, err
	}

	return db, nil
}

func seedDefaultAdmin(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := utils.HashPassword("Admin@123456")
	if err != nil {
		return err
	}

	admin := models.User{
		Username:     "admin",
		PasswordHash: hash,
		DisplayName:  "超级管理员",
		Role:         "admin",
		Status:       1,
	}
	if err := db.Create(&admin).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		return err
	}

	return nil
}

func seedFoundationData(db *gorm.DB) error {
	var roleCount int64
	if err := db.Model(&models.Role{}).Where("role_key = ?", "admin").Count(&roleCount).Error; err != nil {
		return err
	}
	if roleCount == 0 {
		if err := db.Create(&models.Role{
			Name:        "超级管理员",
			RoleKey:     "admin",
			Permissions: "*",
			DataScope:   "all",
			Status:      1,
			Remark:      "系统内置角色",
		}).Error; err != nil {
			return err
		}
	} else {
		if err := db.Model(&models.Role{}).Where("role_key = ? AND permissions = ''", "admin").Update("permissions", "*").Error; err != nil {
			return err
		}
		if err := db.Model(&models.Role{}).Where("role_key = ? AND data_scope = ''", "admin").Update("data_scope", "all").Error; err != nil {
			return err
		}
	}

	var operatorRoleCount int64
	if err := db.Model(&models.Role{}).Where("role_key = ?", "operator").Count(&operatorRoleCount).Error; err != nil {
		return err
	}
	if operatorRoleCount == 0 {
		if err := db.Create(&models.Role{
			Name:        "操作员",
			RoleKey:     "operator",
			Permissions: "user:list,department:list,position:list,notice:list,online-user:list",
			DataScope:   "dept",
			Status:      1,
			Remark:      "系统内置只读角色",
		}).Error; err != nil {
			return err
		}
	} else {
		if err := db.Model(&models.Role{}).Where("role_key = ? AND data_scope = ''", "operator").Update("data_scope", "dept").Error; err != nil {
			return err
		}
	}

	var paramCount int64
	if err := db.Model(&models.SystemParam{}).Where("param_key = ?", "sys.siteName").Count(&paramCount).Error; err != nil {
		return err
	}
	if paramCount == 0 {
		if err := db.Create(&models.SystemParam{
			ParamKey:   "sys.siteName",
			ParamValue: "Tencent Admin Starter",
			ParamName:  "站点名称",
			Status:     1,
			Remark:     "系统默认参数",
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

type schemaMigration struct {
	ID        uint      `gorm:"primaryKey"`
	Filename  string    `gorm:"size:255;uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func runSQLMigrations(db *gorm.DB, dir string) error {
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		var count int64
		if err := db.Model(&schemaMigration{}).Where("filename = ?", name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		contentBytes, readErr := os.ReadFile(filepath.Join(dir, name))
		if readErr != nil {
			return readErr
		}

		stmts := strings.Split(string(contentBytes), ";")
		for _, stmt := range stmts {
			lines := strings.Split(stmt, "\n")
			filtered := make([]string, 0, len(lines))
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if trimmed == "" || strings.HasPrefix(trimmed, "--") {
					continue
				}
				filtered = append(filtered, line)
			}
			sql := strings.TrimSpace(strings.Join(filtered, "\n"))
			if sql == "" {
				continue
			}
			if err := db.Exec(sql).Error; err != nil {
				msg := err.Error()
				if strings.Contains(msg, "already exists") || strings.Contains(msg, "Duplicate entry") {
					continue
				}
				return err
			}
		}

		if err := db.Create(&schemaMigration{Filename: name, CreatedAt: time.Now()}).Error; err != nil {
			return err
		}
	}

	return nil
}
