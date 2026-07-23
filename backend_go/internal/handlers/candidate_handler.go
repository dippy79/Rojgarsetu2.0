package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/middleware"
	"github.com/rojgarsetu/backend/internal/services"
)

type CandidateHandler struct {
	svc *services.CandidateService
}

func NewCandidateHandler(svc *services.CandidateService) *CandidateHandler {
	return &CandidateHandler{svc: svc}
}

func (h *CandidateHandler) ListCandidates(c *gin.Context) {
	location := c.Query("location")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	rows, count, err := h.svc.ListCandidates(c, location, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rows, "count": count})
}

func (h *CandidateHandler) GetCandidate(c *gin.Context) {
	id := c.Param("id")
	candidate, err := h.svc.GetCandidateByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": candidate})
}

func (h *CandidateHandler) GetMyProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	candidate, err := h.svc.GetCandidateByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": candidate})
}

func (h *CandidateHandler) UpdateMyProfile(c *gin.Context) {
	var req db.UpdateCandidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	candidate, err := h.svc.GetCandidateByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}
	updated, err := h.svc.UpdateCandidate(c, candidate.ID.String(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}
