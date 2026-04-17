package models

import "time"

type Role struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:64;not null"`
	RoleKey     string    `json:"roleKey" gorm:"size:64;uniqueIndex;not null"`
	Permissions string    `json:"permissions" gorm:"size:1000;not null;default:''"`
	DataScope   string    `json:"dataScope" gorm:"size:32;not null;default:self"`
	Status      int       `json:"status" gorm:"not null"`
	Remark      string    `json:"remark" gorm:"size:255"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	RoleKey     string `json:"roleKey" binding:"required"`
	Permissions string `json:"permissions"`
	DataScope   string `json:"dataScope"`
	Status      *int   `json:"status"`
	Remark      string `json:"remark"`
}

type UpdateRoleRequest struct {
	Name        string `json:"name"`
	RoleKey     string `json:"roleKey"`
	Permissions string `json:"permissions"`
	DataScope   string `json:"dataScope"`
	Status      *int   `json:"status"`
	Remark      string `json:"remark"`
}
