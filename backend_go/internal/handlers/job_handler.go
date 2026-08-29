package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
)

type JobHandler struct {
	svc *services.JobService
}

func NewJobHandler(svc *services.JobService) *JobHandler {
	return &JobHandler{svc: svc}
}

func (h *JobHandler) ListActiveJobs(c *gin.Context) {
	location := c.Query("location")
	jobType := c.Query("jobType")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	rows, count, err := h.svc.ListActiveJobs(c, location, jobType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "count": count})
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id := c.Param("id")
	job, err := h.svc.GetJobByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	_ = h.svc.IncrementViews(c, id) // Fire and forget
	c.JSON(http.StatusOK, gin.H{"data": job})
}

// CreateJob godoc
// @Summary Create a new job
// @Tags Jobs
// @Security CookieAuth
// @Param request body db.CreateJobRequest true "Job Data"
// @Success 201 {object} db.SavedJob
// @Router /jobs [post]
func (h *JobHandler) CreateJob(c *gin.Context) {
	var req db.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	job, err := h.svc.CreateJob(c, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": job})
}

// UpdateJob godoc
// @Summary Update a job
// @Tags Jobs
// @Security CookieAuth
// @Param id path string true "Job ID"
// @Param request body db.CreateJobRequest true "Job Data"
// @Success 200 {object} db.SavedJob
// @Router /jobs/{id} [put]
func (h *JobHandler) UpdateJob(c *gin.Context) {
	id := c.Param("id")
	var req db.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	job, err := h.svc.UpdateJob(c, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": job})
}

// DeleteJob godoc
// @Summary Delete a job
// @Tags Jobs
// @Security CookieAuth
// @Param id path string true "Job ID"
// @Success 204
// @Router /jobs/{id} [delete]
func (h *JobHandler) DeleteJob(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.DeleteJob(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (h *JobHandler) GetMyJobs(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	// Get company ID from userID via company_service.GetCompanyByUserID
	// Then GetJobsByCompanyID
	// Placeholder - implement full
	c.JSON(http.StatusOK, gin.H{"data": []any{}, "count": 0})
}
