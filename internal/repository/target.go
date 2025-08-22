package repository

import (
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"

	"gorm.io/gorm"
)

type TargetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

func (r *TargetRepository) Create(target *models.Target) error {
	if err := r.db.Create(target).Error; err != nil {
		return custerr.NewInternalErr(err)
	}
	return nil
}

func (r *TargetRepository) GetByID(id uint) (*models.Target, error) {
	var target models.Target
	err := r.db.First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) Update(target *models.Target) error {
	return r.db.Save(target).Error
}

func (r *TargetRepository) Delete(id uint) error {
	return r.db.Delete(&models.Target{}, id).Error
}

func (r *TargetRepository) CountByMissionID(missionID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Target{}).Where("mission_id = ?", missionID).Count(&count).Error
	return count, err
}
