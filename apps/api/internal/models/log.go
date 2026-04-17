package models

import "time"

type OperationLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"userId" gorm:"index"`
	Username   string    `json:"username" gorm:"size:64;index"`
	Method     string    `json:"method" gorm:"size:16;not null"`
	Path       string    `json:"path" gorm:"size:255;not null"`
	StatusCode int       `json:"statusCode" gorm:"not null"`
	LatencyMs  int64     `json:"latencyMs" gorm:"not null"`
	IP         string    `json:"ip" gorm:"size:64"`
	UserAgent  string    `json:"userAgent" gorm:"size:255"`
	CreatedAt  time.Time `json:"createdAt"`
}

type LoginLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"size:64;index"`
	Status    int       `json:"status" gorm:"not null"`
	Message   string    `json:"message" gorm:"size:255"`
	IP        string    `json:"ip" gorm:"size:64"`
	UserAgent string    `json:"userAgent" gorm:"size:255"`
	CreatedAt time.Time `json:"createdAt"`
}
