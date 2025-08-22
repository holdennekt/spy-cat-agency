package tests

import (
	"spy-cat-agency/internal/models"
	"spy-cat-agency/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMissionRepository struct {
	mock.Mock
}

func (m *MockMissionRepository) Create(mission *models.Mission) error {
	args := m.Called(mission)
	return args.Error(0)
}

func (m *MockMissionRepository) GetAll() ([]models.Mission, error) {
	args := m.Called()
	return args.Get(0).([]models.Mission), args.Error(1)
}

func (m *MockMissionRepository) GetByID(id uint) (*models.Mission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Mission), args.Error(1)
}

func (m *MockMissionRepository) Update(mission *models.Mission) error {
	args := m.Called(mission)
	return args.Error(0)
}

func (m *MockMissionRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockMissionRepository) AssignCat(missionID, catID uint) error {
	args := m.Called(missionID, catID)
	return args.Error(0)
}

type MockTargetRepository struct {
	mock.Mock
}

func (m *MockTargetRepository) Create(target *models.Target) error {
	args := m.Called(target)
	return args.Error(0)
}

func (m *MockTargetRepository) GetByID(id uint) (*models.Target, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Target), args.Error(1)
}

func (m *MockTargetRepository) Update(target *models.Target) error {
	args := m.Called(target)
	return args.Error(0)
}

func (m *MockTargetRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTargetRepository) CountByMissionID(missionID uint) (int64, error) {
	args := m.Called(missionID)
	return args.Get(0).(int64), args.Error(1)
}

func TestMissionService_CreateMission(t *testing.T) {
	mockMissionRepo := new(MockMissionRepository)
	mockTargetRepo := new(MockTargetRepository)
	missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

	t.Run("successful creation with valid targets", func(t *testing.T) {
		dto := models.CreateMissionDTO{
			Targets: []models.CreateTargetDTO{
				{Name: "Target 1", Country: "USA", Notes: "Important target"},
				{Name: "Target 2", Country: "UK", Notes: "Secondary target"},
			},
		}

		mission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("Create", mock.AnythingOfType("*models.Mission")).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.Mission)
			arg.ID = 1
		})

		mockTargetRepo.On("Create", mock.AnythingOfType("*models.Target")).Return(nil).Times(2)

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil)

		result, err := missionService.Create(dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("invalid number of targets", func(t *testing.T) {
		dto := models.CreateMissionDTO{
			Targets: []models.CreateTargetDTO{},
		}

		result, err := missionService.Create(dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "mission must have 1-3 targets", err.Error())
	})
}

func TestMissionService_DeleteMission(t *testing.T) {
	mockTargetRepo := new(MockTargetRepository)
	mockMissionRepo := new(MockMissionRepository)
	missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

	t.Run("successful deletion of unassigned mission", func(t *testing.T) {
		mission := &models.Mission{ID: 1, CatID: nil, Complete: false}

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil).Once()
		mockMissionRepo.On("Delete", uint(1)).Return(nil)

		err := missionService.Delete(1)

		assert.NoError(t, err)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cannot delete assigned mission", func(t *testing.T) {
		catID := uint(1)
		mission := &models.Mission{ID: 1, CatID: &catID, Complete: false}

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil).Once()

		err := missionService.Delete(1)

		assert.Error(t, err)
		assert.Equal(t, "cannot delete assigned mission", err.Error())
		mockMissionRepo.AssertExpectations(t)
	})
}
