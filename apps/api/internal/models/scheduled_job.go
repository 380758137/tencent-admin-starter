package models

import "time"

type ScheduledJob struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"size:128;not null"`
	CronExpr   string     `json:"cronExpr" gorm:"size:64;not null"`
	Command    string     `json:"command" gorm:"size:255;not null"`
	Status     int        `json:"status" gorm:"not null;default:1"`
	LastRunAt  *time.Time `json:"lastRunAt"`
	LastResult string     `json:"lastResult" gorm:"size:32"`
	Remark     string     `json:"remark" gorm:"size:255"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type ScheduledJobLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	JobID       uint      `json:"jobId" gorm:"index;not null"`
	TriggerType string    `json:"triggerType" gorm:"size:32;not null"`
	Status      string    `json:"status" gorm:"size:32;not null"`
	Message     string    `json:"message" gorm:"size:255"`
	RunAt       time.Time `json:"runAt" gorm:"not null"`
}

type CreateScheduledJobRequest struct {
	Name     string `json:"name" binding:"required"`
	CronExpr string `json:"cronExpr" binding:"required"`
	Command  string `json:"command" binding:"required"`
	Status   *int   `json:"status"`
	Remark   string `json:"remark"`
}

type UpdateScheduledJobRequest struct {
	Name     string `json:"name"`
	CronExpr string `json:"cronExpr"`
	Command  string `json:"command"`
	Status   *int   `json:"status"`
	Remark   string `json:"remark"`
}
