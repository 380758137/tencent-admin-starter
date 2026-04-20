package models

import "time"

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	DisplayName  string    `json:"displayName" gorm:"size:64;not null"`
	Role         string    `json:"role" gorm:"size:255;not null;default:operator"`
	DeptID       uint      `json:"deptId" gorm:"not null;default:0;index"`
	Status       int       `json:"status" gorm:"not null"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName" binding:"required"`
	Role        string `json:"role"`
	DeptID      *uint  `json:"deptId"`
	Status      *int   `json:"status"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"displayName"`
	Role        string `json:"role"`
	DeptID      *uint  `json:"deptId"`
	Status      *int   `json:"status"`
	Password    string `json:"password"`
}
