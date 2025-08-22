package service

import (
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"
)

type MissionService struct {
	missionRepo MissionRepository
	targetRepo  TargetRepository
}

func NewMissionService(missionRepo MissionRepository, targetRepo TargetRepository) *MissionService {
	return &MissionService{
		missionRepo: missionRepo,
		targetRepo:  targetRepo,
	}
}

func (s *MissionService) Create(dto models.CreateMissionDTO) (*models.Mission, error) {
	mission := &models.Mission{}

	if err := s.missionRepo.Create(mission); err != nil {
		return nil, err
	}

	for _, targetReq := range dto.Targets {
		target := &models.Target{
			MissionID: mission.ID,
			Name:      targetReq.Name,
			Country:   targetReq.Country,
			Notes:     targetReq.Notes,
		}
		if err := s.targetRepo.Create(target); err != nil {
			return nil, err
		}
	}

	return s.missionRepo.GetByID(mission.ID)
}

func (s *MissionService) GetAll() ([]models.Mission, error) {
	return s.missionRepo.GetAll()
}

func (s *MissionService) GetByID(id uint) (*models.Mission, error) {
	return s.missionRepo.GetByID(id)
}

func (s *MissionService) Update(id uint, dto models.UpdateMissionDTO) (*models.Mission, error) {
	mission, err := s.missionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if dto.Complete != nil {
		mission.Complete = *dto.Complete
	}

	return mission, s.missionRepo.Update(mission)
}

func (s *MissionService) Delete(id uint) error {
	mission, err := s.missionRepo.GetByID(id)
	if err != nil {
		return err
	}

	if mission.CatID != nil {
		return custerr.NewConflictErr("cannot delete assigned mission")
	}

	return s.missionRepo.Delete(id)
}

func (s *MissionService) AssignCat(missionID, catID uint) (*models.Mission, error) {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, err
	}

	if mission.Complete {
		return nil, custerr.NewConflictErr("cannot assign cat to completed mission")
	}

	activeMission, err := s.missionRepo.GetActiveByCatID(catID)
	if err != nil {
		return nil, err
	}
	if activeMission != nil {
		return nil, custerr.NewConflictErr("cat already has an active mission")
	}

	mission.CatID = &catID

	return mission, s.missionRepo.Update(mission)
}

func (s *MissionService) CreateTarget(missionID uint, dto models.CreateTargetDTO) (*models.Mission, error) {
	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, err
	}

	if mission.Complete {
		return nil, custerr.NewConflictErr("cannot add targets to completed mission")
	}

	count, err := s.targetRepo.CountByMissionID(missionID)
	if err != nil {
		return nil, err
	}

	if count >= 3 {
		return nil, custerr.NewConflictErr("mission cannot have more than 3 targets")
	}

	target := &models.Target{
		MissionID: missionID,
		Name:      dto.Name,
		Country:   dto.Country,
		Notes:     dto.Notes,
	}

	if err := s.targetRepo.Create(target); err != nil {
		return nil, err
	}

	return s.missionRepo.GetByID(missionID)
}

func (s *MissionService) UpdateTarget(missionID, targetID uint, dto models.UpdateTargetDTO) (*models.Mission, error) {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return nil, err
	}

	if target.MissionID != missionID {
		return nil, custerr.NewBadRequestErr("target does not belong to this mission")
	}

	mission, err := s.missionRepo.GetByID(missionID)
	if err != nil {
		return nil, err
	}

	if target.Complete {
		return nil, custerr.NewConflictErr("cannot update notes of completed target")
	}
	if mission.Complete {
		return nil, custerr.NewConflictErr("cannot update notes if mission is completed")
	}

	if dto.Notes != nil {
		target.Notes = *dto.Notes
	}
	if dto.Complete != nil {
		target.Complete = *dto.Complete
	}

	if err := s.targetRepo.Update(target); err != nil {
		return nil, err
	}

	return s.missionRepo.GetByID(missionID)
}

func (s *MissionService) DeleteTarget(missionID, targetID uint) (*models.Mission, error) {
	target, err := s.targetRepo.GetByID(targetID)
	if err != nil {
		return nil, err
	}

	if target.MissionID != missionID {
		return nil, custerr.NewConflictErr("target does not belong to this mission")
	}

	if target.Complete {
		return nil, custerr.NewConflictErr("cannot delete completed target")
	}

	count, err := s.targetRepo.CountByMissionID(missionID)
	if err != nil {
		return nil, err
	}

	if count <= 1 {
		return nil, custerr.NewConflictErr("mission must have at least 1 target")
	}

	if err := s.targetRepo.Delete(targetID); err != nil {
		return nil, err
	}

	return s.missionRepo.GetByID(missionID)
}
