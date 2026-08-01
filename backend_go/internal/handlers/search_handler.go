package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type SearchHandler struct {
	service *services.SearchService
}

func NewSearchHandler(service *services.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var req services.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Try to parse query from query string if JSON body is empty
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, db.ErrorResponse(400, "Search query is required"))
			return
		}
		req.Query = query
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	req.Page = page
	req.Limit = limit

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, "Search query is required"))
		return
	}

	result, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(result, &db.Pagination{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: (result.Total + result.Limit - 1) / result.Limit,
	}))
}

func (h *SearchHandler) SearchGET(c *gin.Context) {
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	req := services.SearchRequest{
		Query: query,
		Page:  page,
		Limit: limit,
	}

	if req.Query == "" {
		c.JSON(http.StatusBadRequest, db.ErrorResponse(400, "Search query is required"))
		return
	}

	result, err := h.service.Search(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, db.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, db.SuccessResponse(result, &db.Pagination{
		Page:       result.Page,
		Limit:      result.Limit,
		Total:      result.Total,
		TotalPages: (result.Total + result.Limit - 1) / result.Limit,
	}))
}
