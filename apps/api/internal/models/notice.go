package models

import "time"

type Notice struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"size:128;not null"`
	NoticeType string    `json:"noticeType" gorm:"size:32;not null;default:notice"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Pinned     int       `json:"pinned" gorm:"not null;default:0"`
	Status     int       `json:"status" gorm:"not null;default:1"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateNoticeRequest struct {
	Title      string `json:"title" binding:"required"`
	NoticeType string `json:"noticeType"`
	Content    string `json:"content" binding:"required"`
	Pinned     *int   `json:"pinned"`
	Status     *int   `json:"status"`
}

type UpdateNoticeRequest struct {
	Title      string `json:"title"`
	NoticeType string `json:"noticeType"`
	Content    string `json:"content"`
	Pinned     *int   `json:"pinned"`
	Status     *int   `json:"status"`
}
