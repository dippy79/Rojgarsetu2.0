package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

// @Summary Get government jobs
// @Description Retrieve paginated list of government jobs with optional filters
// @Tags jobs
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param department query string false "Filter by department"
// @Param location query string false "Filter by location"
// @Param source query string false "Filter by source"
// @Success 200 {object} db.SuccessResponse
// @Router /api/v1/gov-jobs [get]
func (h *GovJobHandler) GetGovJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.GovJobFilter{
		Department: c.Query("department"),
		Location:   c.Query("location"),
		Source:     c.Query("source"),
	}

	jobs, total, err := h.service.GetGovJobs(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch jobs"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	respJobs := make([]GovJobResponse, 0, len(jobs))
	for _, j := range jobs {
		respJobs = append(respJobs, toGovJobResponse(j))
	}

	c.JSON(http.StatusOK, db.SuccessResponse(respJobs, &pagination))
}

type GovJobHandler struct {
	service *services.GovJobService
}

func NewGovJobHandler(service *services.GovJobService) *GovJobHandler {
	return &GovJobHandler{service: service}
}

func (h *GovJobHandler) GetGovJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.GovJobFilter{
		Department: c.Query("department"),
		Location:   c.Query("location"),
		Source:     c.Query("source"),
	}

	jobs, total, err := h.service.GetGovJobs(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch jobs"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	respJobs := make([]GovJobResponse, 0, len(jobs))
	for _, j := range jobs {
		respJobs = append(respJobs, toGovJobResponse(j))
	}

	c.JSON(http.StatusOK, db.SuccessResponse(respJobs, &pagination))
}

func (h *GovJobHandler) GetGovJobByID(c *gin.Context) {
	id := c.Param("id")

	job, err := h.service.GetGovJobByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(toGovJobByIDResponse(*job), nil))
}
