package models

import (
	"time"

	"gorm.io/gorm"
)

type Cat struct {
	ID uint `json:"id" gorm:"primarykey"`
	CreateCatDTO
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Mission *Mission `json:"mission,omitempty" gorm:"foreignkey:CatID"`
}

type CreateCatDTO struct {
	Name            string  `json:"name" gorm:"not null" binding:"required"`
	YearsExperience int     `json:"years_experience" gorm:"not null" binding:"required,min=0"`
	Breed           string  `json:"breed" gorm:"not null" binding:"required"`
	Salary          float64 `json:"salary" gorm:"not null" binding:"required,min=0"`
}

type UpdateCatDTO struct {
	Salary float64 `json:"salary" binding:"required,min=0"`
}
