package handler

import (
	"net/http"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateMission creates a new mission with targets
// @Summary Create a new mission
// @Description Create a new mission with 1-3 targets
// @Tags Missions
// @Accept json
// @Produce json
// @Param dto body models.CreateMissionDTO true "Mission data"
// @Success 201 {object} models.Mission
// @Failure 400 {object} map[string]string
// @Router /missions [post]
func (h *Handler) CreateMission(c *gin.Context) {
	var dto models.CreateMissionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	mission, err := h.missionService.Create(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// GetMissions retrieves all missions
// @Summary Get all missions
// @Description Get a list of all missions with cats and targets
// @Tags Missions
// @Accept json
// @Produce json
// @Success 200 {array} models.Mission
// @Failure 500 {object} map[string]string
// @Router /missions [get]
func (h *Handler) GetMissions(c *gin.Context) {
	missions, err := h.missionService.GetAll()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, missions)
}

// GetMission retrieves a mission by ID
// @Summary Get mission by ID
// @Description Get a single mission by its ID with cat and targets
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200 {object} models.Mission
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id} [get]
func (h *Handler) GetMission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid id"))
		return
	}

	mission, err := h.missionService.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mission)
}

// UpdateMission updates mission status
// @Summary Update mission status
// @Description Update mission completion status
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param dto body models.UpdateMissionDTO true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id} [patch]
func (h *Handler) UpdateMission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var dto models.UpdateMissionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	mission, err := h.missionService.Update(uint(id), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mission)
}

// DeleteMission removes a mission
// @Summary Delete mission
// @Description Delete a mission (only if not assigned to a cat)
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id} [delete]
func (h *Handler) DeleteMission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid id"))
		return
	}

	if err := h.missionService.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}

// AssignCatToMission assigns a cat to a mission
// @Summary Assign cat to mission
// @Description Assign a spy cat to a mission (cat can only have one active mission)
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param cat_id path int true "Cat ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id}/assign/{cat_id} [patch]
func (h *Handler) AssignCatToMission(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid mission id"))
		return
	}

	catID, err := strconv.ParseUint(c.Param("cat_id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid cat id"))
		return
	}

	mission, err := h.missionService.AssignCat(uint(missionID), uint(catID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mission)
}

// CreateTarget adds a target to a mission
// @Summary Add target to mission
// @Description Add a new target to an existing mission (max 3 targets per mission)
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param dto body models.CreateTargetDTO true "Target data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id}/targets [post]
func (h *Handler) CreateTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid mission id"))
		return
	}

	var dto models.CreateTargetDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	mission, err := h.missionService.CreateTarget(uint(missionID), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, mission)
}

// UpdateTarget updates a target
// @Summary Update target
// @Description Update target notes or completion status (cannot update if completed)
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param target_id path int true "Target ID"
// @Param dto body models.UpdateTargetDTO true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id}/targets/{target_id} [patch]
func (h *Handler) UpdateTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid mission id"))
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid target id"))
		return
	}

	var dto models.UpdateTargetDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	mission, err := h.missionService.UpdateTarget(uint(missionID), uint(targetID), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mission)
}

// DeleteTarget removes a target from mission
// @Summary Delete target
// @Description Delete a target from mission (cannot delete if completed, min 1 target required)
// @Tags Missions
// @Accept json
// @Produce json
// @Param id path int true "Mission ID"
// @Param target_id path int true "Target ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /missions/{id}/targets/{target_id} [delete]
func (h *Handler) DeleteTarget(c *gin.Context) {
	missionID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid mission id"))
		return
	}

	targetID, err := strconv.ParseUint(c.Param("target_id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid target id"))
		return
	}

	mission, err := h.missionService.DeleteTarget(uint(missionID), uint(targetID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, mission)
}
