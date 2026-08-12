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

func (h *JobCategoryHandler) GetJobCategories(c *gin.Context) {
	categories, err := h.service.GetJobCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch job categories"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(categories, nil))
}

func (h *JobCategoryHandler) GetJobCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	category, err := h.service.GetJobCategoryBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job category not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(category, nil))
}

func (h *JobCategoryHandler) GetJobCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.service.GetJobCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job category not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(category, nil))
}
