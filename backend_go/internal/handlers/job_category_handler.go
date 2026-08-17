package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type JobCategoryHandler struct {
	service *services.JobCategoryService
}

func NewJobCategoryHandler(service *services.JobCategoryService) *JobCategoryHandler {
	return &JobCategoryHandler{service: service}
}

// @Summary Get job categories
// @Description Retrieve list of job categories
// @Tags categories
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories [get]
func (h *JobCategoryHandler) GetJobCategories(c *gin.Context) {
	categories, err := h.service.GetJobCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch job categories"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(categories, nil))
}

// @Summary Get job category by slug
// @Description Retrieve a specific job category by slug
// @Tags categories
// @Param slug path string true "Category slug"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/slug/{slug} [get]
func (h *JobCategoryHandler) GetJobCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.service.GetJobCategoryBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job category not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(category, nil))
}

// @Summary Get job category by ID
// @Description Retrieve a specific job category by ID
// @Tags categories
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/categories/{id} [get]
func (h *JobCategoryHandler) GetJobCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.service.GetJobCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job category not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(category, nil))
}
