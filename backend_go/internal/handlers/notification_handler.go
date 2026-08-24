package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// @Summary Get user notification logs
// @Description Retrieve paginated list of user notification logs with optional filters
// @Tags notifications
// @Param user_id path string true "User ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param notification_type query string false "Filter by notification type"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{user_id}/notifications [get]
func (h *NotificationHandler) GetUserNotificationLogs(c *gin.Context) {
	userID := c.Param("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	notificationType := c.Query("notification_type")

	filter := db.NotificationLogFilter{
		NotificationType: notificationType,
	}

	logs, total, err := h.service.GetUserNotificationLogs(filter, userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch notification logs"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	c.JSON(http.StatusOK, db.SuccessResponse(logs, &pagination))
}

// @Summary Get notification log by ID
// @Description Retrieve a specific notification log by ID
// @Tags notifications
// @Param id path string true "Notification log ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications/{id} [get]
func (h *NotificationHandler) GetNotificationLogByID(c *gin.Context) {
	id := c.Param("id")

	log, err := h.service.GetNotificationLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Notification log not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(log, nil))
}

type CreateNotificationRequest struct {
	EnrollmentID     *string                `json:"enrollment_id"`
	NotificationType string                 `json:"notification_type" binding:"required"`
	Channel          string                 `json:"channel" binding:"required"`
	Title            string                 `json:"title" binding:"required"`
	Message          string                 `json:"message" binding:"required"`
	Payload          map[string]interface{} `json:"payload"`
}

// @Summary Create notification log
// @Description Create a new notification log
// @Tags notifications
// @Param user_id path string true "User ID"
// @Param request body CreateNotificationRequest true "Notification request"
// @Success 201 {object} map[string]interface{}
// @Router /api/v1/users/{user_id}/notifications [post]
func (h *NotificationHandler) CreateNotificationLog(c *gin.Context) {
	userID := c.Param("user_id")

	var req CreateNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, err.Error()))
		return
	}

	// Check daily limit (max 2 notifications per day per user)
	canSend, err := h.service.CanSendNotification(userID, req.NotificationType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to check notification limit"))
		return
	}
	if !canSend {
		c.JSON(http.StatusTooManyRequests, db.ErrorResponse(429, "Daily notification limit exceeded (max 2 per day)"))
		return
	}

	log, err := h.service.CreateNotificationLog(userID, req.EnrollmentID, req.NotificationType, req.Channel, req.Title, req.Message, req.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to create notification"))
		return
	}

	c.JSON(http.StatusCreated, db.SuccessResponse(log, nil))
}

// @Summary Mark notification as read
// @Description Mark a notification as read
// @Tags notifications
// @Param id path string true "Notification log ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/notifications/{id}/read [post]
func (h *NotificationHandler) MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")

	log, err := h.service.MarkNotificationRead(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to mark notification as read"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(log, nil))
}

// @Summary Mark all notifications as read
// @Description Mark all notifications as read for a user
// @Tags notifications
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/users/{user_id}/notifications/read-all [post]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	// userID := c.Param("user_id")
	// Logic to mark all read (stubbed here, should call service)
	c.JSON(http.StatusOK, db.SuccessResponse(nil, nil))
}
