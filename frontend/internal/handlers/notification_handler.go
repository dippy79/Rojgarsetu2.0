package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "rojgarsetu/internal/services"
)

type NotificationHandler struct {
    Service *services.NotificationService
}

func (h *NotificationHandler) List(c *gin.Context) {
    userID := c.GetString("user_id")
    notifications, err := h.Service.GetUserNotifications(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
    id := c.Param("id")
    if err := h.Service.MarkAsRead(c.Request.Context(), id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "success"})
}
