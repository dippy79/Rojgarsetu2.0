package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type UserEnrollmentHandler struct {
	service *services.UserEnrollmentService
}

func NewUserEnrollmentHandler(service *services.UserEnrollmentService) *UserEnrollmentHandler {
	return &UserEnrollmentHandler{service: service}
}

func (h *UserEnrollmentHandler) GetUserEnrollments(c *gin.Context) {
	userID := c.Param("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")

	filter := db.UserEnrollmentFilter{
		Status: status,
	}

	enrollments, total, err := h.service.GetUserEnrollments(filter, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch enrollments"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	c.JSON(http.StatusOK, db.SuccessResponse(enrollments, &pagination))
}

func (h *UserEnrollmentHandler) GetUserEnrollmentByID(c *gin.Context) {
	id := c.Param("id")

	enrollment, err := h.service.GetUserEnrollmentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Enrollment not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

type CreateEnrollmentRequest struct {
	TradeID   string                 `json:"trade_id" binding:"required"`
	ExpiresAt time.Time              `json:"expires_at" binding:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (h *UserEnrollmentHandler) CreateUserEnrollment(c *gin.Context) {
	userID := c.Param("user_id")

	var req CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, err.Error()))
		return
	}

	enrollment, err := h.service.CreateUserEnrollment(userID, req.TradeID, req.ExpiresAt, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to create enrollment"))
		return
	}

	c.JSON(http.StatusCreated, db.SuccessResponse(enrollment, nil))
}

type UpdateEnrollmentRequest struct {
	Status      string                 `json:"status"`
	ExpiresAt   *time.Time             `json:"expires_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	ProgressPct int32                  `json:"progress_pct"`
	Metadata    map[string]interface{} `json:"metadata"`
}

func (h *UserEnrollmentHandler) UpdateUserEnrollment(c *gin.Context) {
	id := c.Param("id")

	var req UpdateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, err.Error()))
		return
	}

	enrollment, err := h.service.UpdateUserEnrollment(id, req.Status, req.ExpiresAt, req.CompletedAt, req.ProgressPct, req.Metadata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to update enrollment"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

func (h *UserEnrollmentHandler) UpdateEnrollmentProgress(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ProgressPct int32 `json:"progress_pct" binding:"required,min=0,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, err.Error()))
		return
	}

	enrollment, err := h.service.UpdateEnrollmentProgress(id, req.ProgressPct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to update progress"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

func (h *UserEnrollmentHandler) CompleteEnrollment(c *gin.Context) {
	id := c.Param("id")

	enrollment, err := h.service.CompleteEnrollment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to complete enrollment"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

func (h *UserEnrollmentHandler) CancelEnrollment(c *gin.Context) {
	id := c.Param("id")

	enrollment, err := h.service.CancelEnrollment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to cancel enrollment"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

// GetExpiringEnrollments returns enrollments expiring in <= 7 days
// This is the main endpoint for the notification engine
func (h *UserEnrollmentHandler) GetExpiringEnrollments(c *gin.Context) {
	enrollments, err := h.service.GetExpiringEnrollmentsWithTrade()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch expiring enrollments"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollments, nil))
}
