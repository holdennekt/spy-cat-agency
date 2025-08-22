package handler

import "spy-cat-agency/internal/models"

type CatService interface {
	Create(dto *models.CreateCatDTO) (*models.Cat, error)
	GetAll() ([]models.Cat, error)
	GetByID(id uint) (*models.Cat, error)
	Update(id uint, dto models.UpdateCatDTO) (*models.Cat, error)
	Delete(id uint) error
}

type MissionService interface {
	Create(dto models.CreateMissionDTO) (*models.Mission, error)
	GetAll() ([]models.Mission, error)
	GetByID(id uint) (*models.Mission, error)
	Update(id uint, dto models.UpdateMissionDTO) (*models.Mission, error)
	Delete(id uint) error
	AssignCat(missionID, catID uint) (*models.Mission, error)
	CreateTarget(missionID uint, dto models.CreateTargetDTO) (*models.Mission, error)
	UpdateTarget(missionID, targetID uint, dto models.UpdateTargetDTO) (*models.Mission, error)
	DeleteTarget(missionID, targetID uint) (*models.Mission, error)
}
