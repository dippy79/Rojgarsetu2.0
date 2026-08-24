package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
)

type UploadHandler struct {
	svc *services.UploadService
}

func NewUploadHandler(svc *services.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	fileType := c.DefaultPostForm("file_type", "document")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request"})
		return
	}
	defer file.Close()

	// Check file size (5MB limit)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "File size exceeds 5MB limit"})
		return
	}

	upload, err := h.svc.SaveFile(c.Request.Context(), userID, fileType, header.Filename, header.Size, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": upload})
}
