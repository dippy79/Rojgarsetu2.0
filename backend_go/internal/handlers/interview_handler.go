package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type InterviewHandler struct {
	service *services.InterviewService
}

func NewInterviewHandler(service *services.InterviewService) *InterviewHandler {
	return &InterviewHandler{service: service}
}

func (h *InterviewHandler) ScheduleInterview(c *gin.Context) {
	var req struct {
		ApplicationID string    `json:"application_id" binding:"required"`
		CandidateID   string    `json:"candidate_id" binding:"required"`
		CompanyID     string    `json:"company_id" binding:"required"`
		ScheduledAt   time.Time `json:"scheduled_at" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	interview, err := h.service.ScheduleInterview(c.Request.Context(), req.ApplicationID, req.CandidateID, req.CompanyID, req.ScheduledAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, db.SuccessResponse(interview, nil))
}

func (h *InterviewHandler) GetInterviewByID(c *gin.Context) {
	id := c.Param("id")
	interview, err := h.service.GetInterviewByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Interview not found"})
		return
	}
	c.JSON(http.StatusOK, db.SuccessResponse(interview, nil))
}
