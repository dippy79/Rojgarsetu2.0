package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type JobTradeHandler struct {
	service *services.JobTradeService
}

func NewJobTradeHandler(service *services.JobTradeService) *JobTradeHandler {
	return &JobTradeHandler{service: service}
}

func (h *JobTradeHandler) GetJobTrades(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	categoryID := c.Query("category_id")
	demandLevel := c.Query("demand_level")

	filter := db.JobTradeFilter{
		CategoryID:  categoryID,
		DemandLevel: demandLevel,
	}

	trades, total, err := h.service.GetJobTrades(filter, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch job trades"))
		return
	}

	pagination := db.NewPagination(page, limit, total)

	c.JSON(http.StatusOK, db.SuccessResponse(trades, &pagination))
}

func (h *JobTradeHandler) GetJobTradeByID(c *gin.Context) {
	id := c.Param("id")

	trade, err := h.service.GetJobTradeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job trade not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(trade, nil))
}

func (h *JobTradeHandler) GetJobTradeBySlug(c *gin.Context) {
	slug := c.Param("slug")
	categoryID := c.Param("category_id")

	trade, err := h.service.GetJobTradeBySlug(slug, categoryID)
	if err != nil {
		c.JSON(http.StatusNotFound, db.ErrorResponse(404, "Job trade not found"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(trade, nil))
}

func (h *JobTradeHandler) GetJobTradesByCategory(c *gin.Context) {
	categoryID := c.Param("category_id")

	trades, err := h.service.GetJobTradesByCategory(categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, "Failed to fetch job trades"))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(trades, nil))
}
