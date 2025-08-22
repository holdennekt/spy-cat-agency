package handler

import (
	"spy-cat-agency/internal/service"
)

type Handler struct {
	catService     CatService
	missionService MissionService
}

func New(services *service.Service) *Handler {
	return &Handler{
		catService:     services.Cat,
		missionService: services.Mission,
	}
}
