package service

import "spy-cat-agency/internal/models"

type CatRepository interface {
	Create(cat *models.Cat) error
	GetAll() ([]models.Cat, error)
	GetByID(id uint) (*models.Cat, error)
	Update(cat *models.Cat) error
	Delete(id uint) error
}

type MissionRepository interface {
	Create(mission *models.Mission) error
	GetAll() ([]models.Mission, error)
	GetByID(id uint) (*models.Mission, error)
	Update(mission *models.Mission) error
	Delete(id uint) error
	GetActiveByCatID(catID uint) (*models.Mission, error)
}

type TargetRepository interface {
	Create(target *models.Target) error
	GetByID(id uint) (*models.Target, error)
	Update(target *models.Target) error
	Delete(id uint) error
	CountByMissionID(missionID uint) (int64, error)
}
