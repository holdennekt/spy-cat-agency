package repository

import (
	"fmt"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"

	"gorm.io/gorm"
)

type MissionRepository struct {
	db *gorm.DB
}

func NewMissionRepository(db *gorm.DB) *MissionRepository {
	return &MissionRepository{db: db}
}

func (r *MissionRepository) Create(mission *models.Mission) error {
	if err := r.db.Create(mission).Error; err != nil {
		return custerr.NewInternalErr(err)
	}
	return nil
}

func (r *MissionRepository) GetAll() ([]models.Mission, error) {
	var missions []models.Mission
	if err := r.db.Preload("Cat").Preload("Targets").Find(&missions).Error; err != nil {
		return nil, custerr.NewInternalErr(err)
	}
	return missions, nil
}

func (r *MissionRepository) GetByID(id uint) (*models.Mission, error) {
	var mission models.Mission
	err := r.db.Preload("Cat").Preload("Targets").First(&mission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, custerr.NewNotFoundErr(fmt.Sprintf("no mission with id \"%d\"", id))
		}
		return nil, custerr.NewInternalErr(err)
	}
	return &mission, nil
}

func (r *MissionRepository) Update(mission *models.Mission) error {
	res := r.db.Model(&models.Mission{}).Where("id = ?", mission.ID).Updates(mission)
	if res.Error != nil {
		switch res.Error {
		case gorm.ErrForeignKeyViolated:
			return custerr.NewNotFoundErr(fmt.Sprintf("no cat with id \"%d\"", *mission.CatID))
		default:
			return custerr.NewInternalErr(res.Error)
		}
	}

	if res.RowsAffected == 0 {
		return custerr.NewNotFoundErr(fmt.Sprintf("no mission with id \"%d\"", mission.ID))
	}

	return nil
}

func (r *MissionRepository) Delete(id uint) error {
	res := r.db.Delete(&models.Mission{}, id)
	if res.Error != nil {
		return custerr.NewInternalErr(res.Error)
	}

	if res.RowsAffected == 0 {
		return custerr.NewNotFoundErr(fmt.Sprintf("no cat with id \"%d\"", id))
	}
	return nil
}

func (r *MissionRepository) GetActiveByCatID(catID uint) (*models.Mission, error) {
	var mission models.Mission
	err := r.db.Preload("Cat").Preload("Targets").Where("cat_id = ? AND complete = ?", catID, false).First(&mission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, custerr.NewInternalErr(err)
	}
	return &mission, nil
}
