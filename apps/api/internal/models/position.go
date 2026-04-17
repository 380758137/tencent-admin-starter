package models

import "time"

type Position struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:64;not null"`
	Code      string    `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Sort      int       `json:"sort" gorm:"not null;default:0"`
	Status    int       `json:"status" gorm:"not null;default:1"`
	Remark    string    `json:"remark" gorm:"size:255"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreatePositionRequest struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code" binding:"required"`
	Sort   *int   `json:"sort"`
	Status *int   `json:"status"`
	Remark string `json:"remark"`
}

type UpdatePositionRequest struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Sort   *int   `json:"sort"`
	Status *int   `json:"status"`
	Remark string `json:"remark"`
}
