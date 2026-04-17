package models

import "time"

type DictionaryItem struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	DictType  string    `json:"dictType" gorm:"size:64;index;not null"`
	DictLabel string    `json:"dictLabel" gorm:"size:128;not null"`
	DictValue string    `json:"dictValue" gorm:"size:128;not null"`
	Sort      int       `json:"sort" gorm:"not null;default:0"`
	Status    int       `json:"status" gorm:"not null;default:1"`
	Remark    string    `json:"remark" gorm:"size:255"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateDictionaryItemRequest struct {
	DictType  string `json:"dictType" binding:"required"`
	DictLabel string `json:"dictLabel" binding:"required"`
	DictValue string `json:"dictValue" binding:"required"`
	Sort      *int   `json:"sort"`
	Status    *int   `json:"status"`
	Remark    string `json:"remark"`
}

type UpdateDictionaryItemRequest struct {
	DictType  string `json:"dictType"`
	DictLabel string `json:"dictLabel"`
	DictValue string `json:"dictValue"`
	Sort      *int   `json:"sort"`
	Status    *int   `json:"status"`
	Remark    string `json:"remark"`
}
