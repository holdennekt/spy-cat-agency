package tests

import (
	"errors"
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

func (m *MockMissionRepository) GetActiveByCatID(catID uint) (*models.Mission, error) {
	args := m.Called(catID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Mission), args.Error(1)
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
	t.Run("successful creation with targets", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

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

		mockTargetRepo.On("Create", mock.AnythingOfType("*models.Target")).Return(nil)

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil)

		result, err := missionService.Create(dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("successful creation with no targets", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		dto := models.CreateMissionDTO{
			Targets: []models.CreateTargetDTO{},
		}

		mission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("Create", mock.AnythingOfType("*models.Mission")).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.Mission)
			arg.ID = 1
		})

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil)

		result, err := missionService.Create(dto)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})
}

func TestMissionService_DeleteMission(t *testing.T) {
	t.Run("successful deletion of unassigned mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		mission := &models.Mission{ID: 1, CatID: nil, Complete: false}

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil)
		mockMissionRepo.On("Delete", uint(1)).Return(nil)

		err := missionService.Delete(1)

		assert.NoError(t, err)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cannot delete assigned mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		catID := uint(1)
		mission := &models.Mission{ID: 1, CatID: &catID, Complete: false}

		mockMissionRepo.On("GetByID", uint(1)).Return(mission, nil)

		err := missionService.Delete(1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete assigned mission")
		mockMissionRepo.AssertExpectations(t)
	})
}

func TestMissionService_AssignCat(t *testing.T) {
	t.Run("successful cat assignment", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(1), uint(1)
		mission := &models.Mission{ID: 1, Complete: false}
		updatedMission := &models.Mission{ID: 1, CatID: &catID, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockMissionRepo.On("GetActiveByCatID", catID).Return(nil, nil)
		mockMissionRepo.On("Update", mock.MatchedBy(func(m *models.Mission) bool {
			return m.CatID != nil && *m.CatID == catID && m.ID == missionID
		})).Return(nil)

		result, err := missionService.AssignCat(missionID, catID)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cat already has active mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(1), uint(1)
		mission := &models.Mission{ID: 1, Complete: false}
		activeMission := &models.Mission{ID: 2, CatID: &catID, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockMissionRepo.On("GetActiveByCatID", catID).Return(activeMission, nil)

		result, err := missionService.AssignCat(missionID, catID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cat already has an active mission")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cannot assign cat to completed mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(1), uint(1)
		mission := &models.Mission{ID: 1, Complete: true}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)

		result, err := missionService.AssignCat(missionID, catID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot assign cat to completed mission")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("mission not found", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(999), uint(1)

		mockMissionRepo.On("GetByID", missionID).Return(nil, errors.New("mission not found"))

		result, err := missionService.AssignCat(missionID, catID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "mission not found")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("database error when checking active mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(1), uint(1)
		mission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockMissionRepo.On("GetActiveByCatID", catID).Return(nil, errors.New("database error"))

		result, err := missionService.AssignCat(missionID, catID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("successful assignment when cat has no active mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, catID := uint(1), uint(1)
		mission := &models.Mission{ID: 1, Complete: false}
		updatedMission := &models.Mission{ID: 1, CatID: &catID, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockMissionRepo.On("GetActiveByCatID", catID).Return(nil, nil)
		mockMissionRepo.On("Update", mock.MatchedBy(func(m *models.Mission) bool {
			return m.CatID != nil && *m.CatID == catID && m.ID == missionID
		})).Return(nil)

		result, err := missionService.AssignCat(missionID, catID)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockMissionRepo.AssertExpectations(t)
	})
}

func TestMissionService_CreateTarget(t *testing.T) {
	t.Run("successful target creation", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID := uint(1)
		dto := models.CreateTargetDTO{
			Name:    "Target Alpha",
			Country: "USA",
			Notes:   "High priority target",
		}

		mission := &models.Mission{ID: 1, Complete: false}
		updatedMission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(1), nil)
		mockTargetRepo.On("Create", mock.MatchedBy(func(target *models.Target) bool {
			return target.MissionID == missionID &&
				target.Name == dto.Name &&
				target.Country == dto.Country &&
				target.Notes == dto.Notes
		})).Return(nil)
		mockMissionRepo.On("GetByID", missionID).Return(updatedMission, nil)

		result, err := missionService.CreateTarget(missionID, dto)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("cannot create target for completed mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID := uint(1)
		dto := models.CreateTargetDTO{Name: "Target Alpha", Country: "USA"}

		mission := &models.Mission{ID: 1, Complete: true}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)

		result, err := missionService.CreateTarget(missionID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot add targets to completed mission")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cannot create more than 3 targets", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID := uint(1)
		dto := models.CreateTargetDTO{Name: "Target Delta", Country: "Canada"}

		mission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(3), nil)

		result, err := missionService.CreateTarget(missionID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "mission cannot have more than 3 targets")
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("mission not found", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID := uint(999)
		dto := models.CreateTargetDTO{Name: "Target Alpha", Country: "USA"}

		mockMissionRepo.On("GetByID", missionID).Return(nil, errors.New("mission not found"))

		result, err := missionService.CreateTarget(missionID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "mission not found")
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("target creation fails", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID := uint(1)
		dto := models.CreateTargetDTO{Name: "Target Alpha", Country: "USA"}

		mission := &models.Mission{ID: 1, Complete: false}

		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(1), nil)
		mockTargetRepo.On("Create", mock.AnythingOfType("*models.Target")).Return(errors.New("database error"))

		result, err := missionService.CreateTarget(missionID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockMissionRepo.AssertExpectations(t)
		mockTargetRepo.AssertExpectations(t)
	})
}

func TestMissionService_UpdateTarget(t *testing.T) {
	t.Run("successful target update - notes only", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)
		newNotes := "Updated notes"
		dto := models.UpdateTargetDTO{Notes: &newNotes}

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}
		mission := &models.Mission{ID: 1, Complete: false}
		updatedMission := &models.Mission{ID: 1, Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockTargetRepo.On("Update", mock.MatchedBy(func(t *models.Target) bool {
			return t.Notes == newNotes && t.ID == targetID
		})).Return(nil)
		mockMissionRepo.On("GetByID", missionID).Return(updatedMission, nil)

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockTargetRepo.AssertExpectations(t)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("successful target update - complete status", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)
		complete := true
		dto := models.UpdateTargetDTO{Complete: &complete}

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}
		mission := &models.Mission{ID: 1, Complete: false}
		updatedMission := &models.Mission{ID: 1, Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)
		mockTargetRepo.On("Update", mock.MatchedBy(func(t *models.Target) bool {
			return t.Complete == complete && t.ID == targetID
		})).Return(nil)
		mockMissionRepo.On("GetByID", missionID).Return(updatedMission, nil)

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockTargetRepo.AssertExpectations(t)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("target does not belong to mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)
		dto := models.UpdateTargetDTO{Notes: new(string)}

		target := &models.Target{ID: 1, MissionID: 2, Name: "Target Alpha", Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "target does not belong to this mission")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("cannot update completed target", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)
		newNotes := "Updated notes"
		dto := models.UpdateTargetDTO{Notes: &newNotes}

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: true}
		mission := &models.Mission{ID: 1, Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot update notes of completed target")
		mockTargetRepo.AssertExpectations(t)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("cannot update target on completed mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)
		newNotes := "Updated notes"
		dto := models.UpdateTargetDTO{Notes: &newNotes}

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}
		mission := &models.Mission{ID: 1, Complete: true}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockMissionRepo.On("GetByID", missionID).Return(mission, nil)

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot update notes if mission is completed")
		mockTargetRepo.AssertExpectations(t)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("target not found", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(999)
		dto := models.UpdateTargetDTO{Notes: new(string)}

		mockTargetRepo.On("GetByID", targetID).Return(nil, errors.New("target not found"))

		result, err := missionService.UpdateTarget(missionID, targetID, dto)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "target not found")
		mockTargetRepo.AssertExpectations(t)
	})
}

func TestMissionService_DeleteTarget(t *testing.T) {
	t.Run("successful target deletion", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}
		updatedMission := &models.Mission{ID: 1, Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(2), nil)
		mockTargetRepo.On("Delete", targetID).Return(nil)
		mockMissionRepo.On("GetByID", missionID).Return(updatedMission, nil)

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.NoError(t, err)
		assert.Equal(t, updatedMission, result)
		mockTargetRepo.AssertExpectations(t)
		mockMissionRepo.AssertExpectations(t)
	})

	t.Run("target does not belong to mission", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 2, Name: "Target Alpha", Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "target does not belong to this mission")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("cannot delete completed target", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: true}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "cannot delete completed target")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("cannot delete last target", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(1), nil)

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "mission must have at least 1 target")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("target not found", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(999)

		mockTargetRepo.On("GetByID", targetID).Return(nil, errors.New("target not found"))

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "target not found")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("database error during count", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(0), errors.New("database error"))

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "database error")
		mockTargetRepo.AssertExpectations(t)
	})

	t.Run("database error during deletion", func(t *testing.T) {
		mockMissionRepo := new(MockMissionRepository)
		mockTargetRepo := new(MockTargetRepository)
		missionService := service.NewMissionService(mockMissionRepo, mockTargetRepo)

		missionID, targetID := uint(1), uint(1)

		target := &models.Target{ID: 1, MissionID: 1, Name: "Target Alpha", Complete: false}

		mockTargetRepo.On("GetByID", targetID).Return(target, nil)
		mockTargetRepo.On("CountByMissionID", missionID).Return(int64(2), nil)
		mockTargetRepo.On("Delete", targetID).Return(errors.New("delete failed"))

		result, err := missionService.DeleteTarget(missionID, targetID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "delete failed")
		mockTargetRepo.AssertExpectations(t)
	})
}
