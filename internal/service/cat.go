package service

import (
	"errors"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/catapi"
)

var ErrInvalidCatBreed = errors.New("invalid cat breed")

type CatService struct {
	repo         CatRepository
	catValidator catapi.CatValidator
}

func NewCatService(repo CatRepository, catValidator catapi.CatValidator) *CatService {
	return &CatService{
		repo:         repo,
		catValidator: catValidator,
	}
}

func (s *CatService) Create(catDTO *models.CreateCatDTO) (*models.Cat, error) {
	isValid, err := s.catValidator.ValidateBreed(catDTO.Breed)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, ErrInvalidCatBreed
	}

	cat := &models.Cat{
		CreateCatDTO: *catDTO,
	}

	if err := s.repo.Create(cat); err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CatService) GetAll() ([]models.Cat, error) {
	return s.repo.GetAll()
}

func (s *CatService) GetByID(id uint) (*models.Cat, error) {
	return s.repo.GetByID(id)
}

func (s *CatService) Update(id uint, dto models.UpdateCatDTO) (*models.Cat, error) {
	cat, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	cat.Salary = dto.Salary

	return cat, s.repo.Update(cat)
}

func (s *CatService) Delete(id uint) error {
	return s.repo.Delete(id)
}
