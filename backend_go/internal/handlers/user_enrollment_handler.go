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

// @Summary Get user enrollments
// @Description Retrieve paginated list of user enrollments with optional filters
// @Tags enrollments
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{user_id}/enrollments [get]
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

// @Summary Get user enrollment by ID
// @Description Retrieve a specific user enrollment by ID
// @Tags enrollments
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id} [get]
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

// @Summary Create user enrollment
// @Description Create a new user enrollment
// @Tags enrollments
// @Param user_id path string true "User ID"
// @Param request body CreateEnrollmentRequest true "Enrollment request"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/users/{user_id}/enrollments [post]
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

// @Summary Update user enrollment
// @Description Update an existing user enrollment
// @Tags enrollments
// @Param id path string true "Enrollment ID"
// @Param request body UpdateEnrollmentRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id} [put]
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

// @Summary Update enrollment progress
// @Description Update the progress percentage of an enrollment
// @Tags enrollments
// @Param id path string true "Enrollment ID"
// @Param request body object true "Progress request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id}/progress [put]
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

// @Summary Complete enrollment
// @Description Mark an enrollment as completed
// @Tags enrollments
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id}/complete [post]
func (h *UserEnrollmentHandler) CompleteEnrollment(c *gin.Context) {
	id := c.Param("id")

	enrollment, err := h.service.CompleteEnrollment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to complete enrollment"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

// @Summary Cancel enrollment
// @Description Cancel an enrollment
// @Tags enrollments
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/{id}/cancel [post]
func (h *UserEnrollmentHandler) CancelEnrollment(c *gin.Context) {
	id := c.Param("id")

	enrollment, err := h.service.CancelEnrollment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to cancel enrollment"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollment, nil))
}

// @Summary Get expiring enrollments
// @Description Retrieve enrollments expiring in <= 7 days
// @Tags enrollments
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/enrollments/expiring [get]
func (h *UserEnrollmentHandler) GetExpiringEnrollments(c *gin.Context) {
	enrollments, err := h.service.GetExpiringEnrollmentsWithTrade()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch expiring enrollments"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(enrollments, nil))
}
