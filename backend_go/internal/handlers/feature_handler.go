package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
)

type FeatureHandler struct {
	Repo *db.FeatureRepository
}

func NewFeatureHandler(repo *db.FeatureRepository) *FeatureHandler {
	return &FeatureHandler{Repo: repo}
}

// CreateReviewHandler - Company review submit karne ke liye
func (h *FeatureHandler) CreateReviewHandler(c *gin.Context) {
	var req db.CompanyReview
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	if err := h.Repo.InsertCompanyReview(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Review submitted successfully",
		"data":    req,
	})
}

// ReportJobHandler - Job flag/report karne ke liye
func (h *FeatureHandler) ReportJobHandler(c *gin.Context) {
	var req db.JobReport
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	if err := h.Repo.InsertJobReport(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit job report: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job report submitted successfully",
		"data":    req,
	})
}

// InternalRatingHandler - Recruiter rating save karne ke liye
func (h *FeatureHandler) InternalRatingHandler(c *gin.Context) {
	var req db.CandidateInternalRating
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload: " + err.Error()})
		return
	}

	if err := h.Repo.InsertCandidateInternalRating(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save internal rating: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Internal candidate rating saved",
		"data":    req,
	})
}
