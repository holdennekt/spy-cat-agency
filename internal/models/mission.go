package models

import (
	"time"

	"gorm.io/gorm"
)

type Mission struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CatID     *uint          `json:"cat_id" gorm:"index"`
	Complete  bool           `json:"complete" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Cat     *Cat     `json:"cat,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Targets []Target `json:"targets" gorm:"foreignkey:MissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type CreateMissionDTO struct {
	Targets []CreateTargetDTO `json:"targets" binding:"required,min=1,max=3"`
}

type UpdateMissionDTO struct {
	Complete *bool `json:"complete"`
}
