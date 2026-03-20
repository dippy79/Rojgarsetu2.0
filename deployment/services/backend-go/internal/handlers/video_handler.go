package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type VideoHandler struct {
	service *services.VideoService
}

func NewVideoHandler(service *services.VideoService) *VideoHandler {
	return &VideoHandler{service: service}
}

func (h *VideoHandler) GetVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.VideoFilter{
		Channel:  c.Query("channel"),
		Category: c.Query("category"),
	}

	videos, total, err := h.service.GetVideos(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch videos"))
		return
	}

	pagination := db.NewPagination(page, limit, total)
	c.JSON(http.StatusOK, db.SuccessResponse(videos, &pagination))
}

func (h *VideoHandler) GetVideoByID(c *gin.Context) {
	id := c.Param("id")

	video, err := h.service.GetVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Video not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(video, nil))
}
