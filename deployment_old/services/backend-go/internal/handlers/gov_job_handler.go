package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type GovJobHandler struct {
	service *services.GovJobService
}

func NewGovJobHandler(service *services.GovJobService) *GovJobHandler {
	return &GovJobHandler{service: service}
}

func (h *GovJobHandler) GetGovJobs(c *gin.Context) {
	// Parse pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Parse filters
	filter := db.GovJobFilter{
		Department: c.Query("department"),
		Location:   c.Query("location"),
		Source:     c.Query("source"),
	}

	// Get jobs
	jobs, total, err := h.service.GetGovJobs(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch jobs"))
		return
	}

	// Return response
	pagination := db.NewPagination(page, limit, total)
	c.JSON(http.StatusOK, db.SuccessResponse(jobs, &pagination))
}

func (h *GovJobHandler) GetGovJobByID(c *gin.Context) {
	id := c.Param("id")

	job, err := h.service.GetGovJobByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(job, nil))
}
