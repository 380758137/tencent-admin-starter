package models

import "time"

type Menu struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ParentID  uint      `json:"parentId" gorm:"not null;default:0;index"`
	Name      string    `json:"name" gorm:"size:64;not null"`
	MenuType  string    `json:"menuType" gorm:"size:16;not null;default:menu"`
	Path      string    `json:"path" gorm:"size:255"`
	Component string    `json:"component" gorm:"size:255"`
	Perms     string    `json:"perms" gorm:"size:128"`
	Sort      int       `json:"sort" gorm:"not null;default:0"`
	Status    int       `json:"status" gorm:"not null;default:1"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateMenuRequest struct {
	ParentID  *uint  `json:"parentId"`
	Name      string `json:"name" binding:"required"`
	MenuType  string `json:"menuType"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Perms     string `json:"perms"`
	Sort      *int   `json:"sort"`
	Status    *int   `json:"status"`
}

type UpdateMenuRequest struct {
	ParentID  *uint   `json:"parentId"`
	Name      *string `json:"name"`
	MenuType  *string `json:"menuType"`
	Path      *string `json:"path"`
	Component *string `json:"component"`
	Perms     *string `json:"perms"`
	Sort      *int    `json:"sort"`
	Status    *int    `json:"status"`
}
