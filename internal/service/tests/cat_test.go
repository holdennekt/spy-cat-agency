package tests

import (
	"spy-cat-agency/internal/models"
	"spy-cat-agency/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCatRepository struct {
	mock.Mock
}

func (m *MockCatRepository) Create(cat *models.Cat) error {
	args := m.Called(cat)
	return args.Error(0)
}

func (m *MockCatRepository) GetAll() ([]models.Cat, error) {
	args := m.Called()
	return args.Get(0).([]models.Cat), args.Error(1)
}

func (m *MockCatRepository) GetByID(id uint) (*models.Cat, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cat), args.Error(1)
}

func (m *MockCatRepository) Update(cat *models.Cat) error {
	args := m.Called(cat)
	return args.Error(0)
}

func (m *MockCatRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockCatValidator struct {
	mock.Mock
}

func (m *MockCatValidator) ValidateBreed(breed string) (bool, error) {
	args := m.Called(breed)
	return args.Bool(0), args.Error(1)
}

func TestCatService_CreateCat(t *testing.T) {
	mockRepo := new(MockCatRepository)
	mockValidator := new(MockCatValidator)
	catService := service.NewCatService(mockRepo, mockValidator)

	catDTO := &models.CreateCatDTO{
		Name:            "Agent Whiskers",
		YearsExperience: 5,
		Breed:           "Persian",
		Salary:          50000,
	}

	t.Run("successful creation", func(t *testing.T) {
		mockValidator.On("ValidateBreed", "Persian").Return(true, nil)
		mockRepo.On("Create", catDTO).Return(nil)

		_, err := catService.Create(catDTO)

		assert.NoError(t, err)
		mockValidator.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid breed", func(t *testing.T) {
		mockValidator.On("ValidateBreed", "InvalidBreed").Return(false, nil)

		invalidCatDTO := &models.CreateCatDTO{
			Name:            "Agent Invalid",
			YearsExperience: 3,
			Breed:           "InvalidBreed",
			Salary:          45000,
		}

		_, err := catService.Create(invalidCatDTO)

		assert.Error(t, err)
		assert.Equal(t, "invalid cat breed", err.Error())
		mockValidator.AssertExpectations(t)
	})
}
