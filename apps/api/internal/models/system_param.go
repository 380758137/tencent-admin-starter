package models

import "time"

type SystemParam struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ParamKey   string    `json:"paramKey" gorm:"size:128;uniqueIndex;not null"`
	ParamValue string    `json:"paramValue" gorm:"size:1000;not null"`
	ParamName  string    `json:"paramName" gorm:"size:128;not null"`
	Status     int       `json:"status" gorm:"not null;default:1"`
	Remark     string    `json:"remark" gorm:"size:255"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateSystemParamRequest struct {
	ParamKey   string `json:"paramKey" binding:"required"`
	ParamValue string `json:"paramValue" binding:"required"`
	ParamName  string `json:"paramName" binding:"required"`
	Status     *int   `json:"status"`
	Remark     string `json:"remark"`
}

type UpdateSystemParamRequest struct {
	ParamValue string `json:"paramValue"`
	ParamName  string `json:"paramName"`
	Status     *int   `json:"status"`
	Remark     string `json:"remark"`
}
