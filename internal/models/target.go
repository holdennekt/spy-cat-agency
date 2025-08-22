package models

import (
	"time"

	"gorm.io/gorm"
)

type Target struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	MissionID uint           `json:"mission_id" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"not null" binding:"required"`
	Country   string         `json:"country" gorm:"not null" binding:"required"`
	Notes     string         `json:"notes"`
	Complete  bool           `json:"complete" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Mission Mission `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CreateTargetDTO struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
	Notes   string `json:"notes"`
}

type UpdateTargetDTO struct {
	Notes    *string `json:"notes"`
	Complete *bool   `json:"complete"`
}
