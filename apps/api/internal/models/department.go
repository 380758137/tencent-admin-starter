package models

import "time"

type Department struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:64;not null"`
	Code      string    `json:"code" gorm:"size:64;uniqueIndex;not null"`
	ParentID  uint      `json:"parentId" gorm:"not null;default:0;index"`
	Manager   string    `json:"manager" gorm:"size:64"`
	Status    int       `json:"status" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateDepartmentRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	ParentID *uint  `json:"parentId"`
	Manager string `json:"manager"`
	Status  *int   `json:"status"`
}

type UpdateDepartmentRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	ParentID *uint  `json:"parentId"`
	Manager string `json:"manager"`
	Status  *int   `json:"status"`
}
