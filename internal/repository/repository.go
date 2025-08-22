package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	Cat     *CatRepository
	Mission *MissionRepository
	Target  *TargetRepository
}

func New(db *gorm.DB) *Repository {
	return &Repository{
		Cat:     NewCatRepository(db),
		Mission: NewMissionRepository(db),
		Target:  NewTargetRepository(db),
	}
}
