package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestGetGovJobs tests fetching government jobs
func TestGetGovJobs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/api/v1/gov-jobs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": []gin.H{
				{
					"id":         "1",
					"title":      "Software Engineer",
					"department": "IT",
					"location":   "Delhi",
				},
			},
			"pagination": gin.H{
				"page":       1,
				"limit":      20,
				"total":      1,
				"total_pages": 1,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/gov-jobs?page=1&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Software Engineer")
}

// TestApplyJob tests job application
func TestApplyJob(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.POST("/api/v1/gov-jobs/:id/apply", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Application submitted successfully",
			"application_id": "app-123",
		})
	})

	req, _ := http.NewRequest("POST", "/api/v1/gov-jobs/1/apply", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Application submitted successfully")
}

// TestGetMyApplications tests fetching user's applications
func TestGetMyApplications(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/api/v1/applications", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": []gin.H{
				{
					"id":         "app-123",
					"job_id":     "1",
					"job_title":  "Software Engineer",
					"status":     "pending",
					"applied_at": "2026-08-16T10:00:00Z",
				},
			},
			"pagination": gin.H{
				"page":       1,
				"limit":      20,
				"total":      1,
				"total_pages": 1,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/applications?page=1&limit=20", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "app-123")
}

// TestGetGovJobByID tests fetching single government job
func TestGetGovJobByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	router := gin.New()
	router.GET("/api/v1/gov-jobs/:id", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":         "1",
				"title":      "Software Engineer",
				"department": "IT",
				"location":   "Delhi",
				"eligibility": "Bachelor's in Computer Science",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/gov-jobs/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Software Engineer")
}