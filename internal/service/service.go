package service

import (
	"spy-cat-agency/config"
	"spy-cat-agency/internal/repository"
	"spy-cat-agency/pkg/catapi"
)

type Service struct {
	Cat     *CatService
	Mission *MissionService
}

func New(repo *repository.Repository, cfg *config.Config) *Service {
	catValidator := catapi.NewCatValidator()

	return &Service{
		Cat:     NewCatService(repo.Cat, catValidator),
		Mission: NewMissionService(repo.Mission, repo.Target),
	}
}
