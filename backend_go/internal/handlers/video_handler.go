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

// @Summary Get videos
// @Description Retrieve paginated list of videos with optional filters
// @Tags videos
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param channel query string false "Filter by channel"
// @Param category query string false "Filter by category"
// @Param exclude query string false "Exclude specific video IDs"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/videos [get]
func (h *VideoHandler) GetVideos(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := db.VideoFilter{
		Channel:  c.Query("channel"),
		Category: c.Query("category"),
	}
	exclude := c.Query("exclude")

	videos, total, err := h.service.GetVideos(filter, exclude, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch videos"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	respVideos := make([]VideoResponse, 0, len(videos))
	for _, v := range videos {
		respVideos = append(respVideos, toVideoResponse(v))
	}

	c.JSON(http.StatusOK, db.SuccessResponse(respVideos, &pagination))
}

// @Summary Get video channels
// @Description Retrieve list of video channels
// @Tags videos
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/videos/channels [get]
func (h *VideoHandler) GetVideoChannels(c *gin.Context) {
	channels, err := h.service.GetVideoChannels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch video channels"))
		return
	}
	c.JSON(http.StatusOK, db.SuccessResponse(channels, nil))
}

// @Summary Get video categories
// @Description Retrieve list of video categories
// @Tags videos
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/videos/categories [get]
func (h *VideoHandler) GetVideoCategories(c *gin.Context) {
	categories, err := h.service.GetVideoCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch video categories"))
		return
	}
	c.JSON(http.StatusOK, db.SuccessResponse(categories, nil))
}

func (h *VideoHandler) GetVideoByID(c *gin.Context) {
	id := c.Param("id")

	video, err := h.service.GetVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Video not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(toVideoByIDResponse(*video), nil))
}
