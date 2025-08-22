package handler

import (
	"net/http"
	"spy-cat-agency/internal/models"
	"spy-cat-agency/pkg/custerr"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateCat creates a new cat
// @Summary Create a new cat
// @Description Create a new cat with breed validation using TheCatAPI
// @Tags Cats
// @Accept json
// @Produce json
// @Param dto body models.CreateCatDTO true "Cat data"
// @Success 201 {object} models.Cat
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cats [post]
func (h *Handler) CreateCat(c *gin.Context) {
	var dto models.CreateCatDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	cat, err := h.catService.Create(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// GetCats retrieves all cats
// @Summary Get all cats
// @Description Get a list of all cats with their missions
// @Tags Cats
// @Accept json
// @Produce json
// @Success 200 {array} models.Cat
// @Failure 500 {object} map[string]string
// @Router /cats [get]
func (h *Handler) GetCats(c *gin.Context) {
	cats, err := h.catService.GetAll()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cats)
}

// GetCat retrieves a cat by ID
// @Summary Get cat by ID
// @Description Get a single cat by its ID with mission details
// @Tags Cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} models.Cat
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cats/{id} [get]
func (h *Handler) GetCat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid id"))
		return
	}

	cat, err := h.catService.GetByID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cat)
}

// UpdateCat updates a cat's salary
// @Summary Update cat salary
// @Description Update a cat's salary (only salary can be modified)
// @Tags Cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Param dto body models.UpdateCatDTO true "Update data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cats/{id} [patch]
func (h *Handler) UpdateCat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid id"))
		return
	}

	var dto models.UpdateCatDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.Error(err).SetType(gin.ErrorTypeBind)
		return
	}

	cat, err := h.catService.Update(uint(id), dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, cat)
}

// DeleteCat removes a cat
// @Summary Delete cat
// @Description Delete a cat (only if not on active mission)
// @Tags Cats
// @Accept json
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cats/{id} [delete]
func (h *Handler) DeleteCat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(custerr.NewBadRequestErr("invalid id"))
		return
	}

	if err := h.catService.Delete(uint(id)); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
