package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type PrivJobHandler struct {
	service *services.PrivJobService
}

func NewPrivJobHandler(service *services.PrivJobService) *PrivJobHandler {
	return &PrivJobHandler{service: service}
}

// @Summary Get private jobs
// @Description Retrieve paginated list of private jobs with optional filters
// @Tags jobs
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param company query string false "Filter by company"
// @Param location query string false "Filter by location"
// @Param source query string false "Filter by source"
// @Param job_type query string false "Filter by job type"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/priv-jobs [get]
func (h *PrivJobHandler) GetPrivJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.PrivJobFilter{
		Company:  c.Query("company"),
		Location: c.Query("location"),
		Source:   c.Query("source"),
		JobType:  c.Query("job_type"),
	}

	jobs, total, err := h.service.GetPrivJobs(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch jobs"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	respJobs := make([]PrivJobResponse, 0, len(jobs))
	for _, j := range jobs {
		respJobs = append(respJobs, toPrivJobResponse(j))
	}

	c.JSON(http.StatusOK, db.SuccessResponse(respJobs, &pagination))
}

func (h *PrivJobHandler) GetPrivJobByID(c *gin.Context) {
	id := c.Param("id")

	job, err := h.service.GetPrivJobByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(toPrivJobByIDResponse(*job), nil))
}
