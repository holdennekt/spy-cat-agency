package repository

import (
	"fmt"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"

	"gorm.io/gorm"
)

type CatRepository struct {
	db *gorm.DB
}

func NewCatRepository(db *gorm.DB) *CatRepository {
	return &CatRepository{db: db}
}

func (r *CatRepository) Create(cat *models.Cat) error {
	if err := r.db.Create(cat).Error; err != nil {
		return custerr.NewInternalErr(err)
	}
	return nil
}

func (r *CatRepository) GetAll() ([]models.Cat, error) {
	var cats []models.Cat
	if err := r.db.Preload("Mission.Targets").Find(&cats).Error; err != nil {
		return nil, custerr.NewInternalErr(err)
	}
	return cats, nil
}

func (r *CatRepository) GetByID(id uint) (*models.Cat, error) {
	var cat models.Cat
	err := r.db.Preload("Mission.Targets").First(&cat, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, custerr.NewNotFoundErr(fmt.Sprintf("no cat with id \"%d\"", id))
		}
		return nil, custerr.NewInternalErr(err)
	}
	return &cat, nil
}

func (r *CatRepository) Update(cat *models.Cat) error {
	if err := r.db.Save(cat).Error; err != nil {
		return custerr.NewInternalErr(err)
	}
	return nil
}

func (r *CatRepository) Delete(id uint) error {
	res := r.db.Delete(&models.Cat{}, id)
	if res.Error != nil {
		return custerr.NewInternalErr(res.Error)
	}

	if res.RowsAffected == 0 {
		return custerr.NewNotFoundErr(fmt.Sprintf("no cat with id \"%d\"", id))
	}
	return nil
}
