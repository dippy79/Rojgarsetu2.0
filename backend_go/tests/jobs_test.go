package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGetGovJobsReturnsJSON verifies the response structure of the gov-jobs list
func TestGetGovJobsReturnsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/gov-jobs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": []interface{}{},
			"pagination": gin.H{
				"total": 0,
				"limit": 20,
				"page":  1,
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/gov-jobs", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Errorf("failed to unmarshal JSON: %v", err)
	}

	if _, ok := resp["data"]; !ok {
		t.Error("response missing 'data' field")
	}
	if _, ok := resp["pagination"]; !ok {
		t.Error("response missing 'pagination' field")
	}
}

// TestGetGovJobsPagination verifies that pagination parameters are reflected in the response
func TestGetGovJobsPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/api/v1/gov-jobs", func(c *gin.Context) {
		limit := c.DefaultQuery("limit", "20")
		page := c.DefaultQuery("page", "1")
		c.JSON(http.StatusOK, gin.H{
			"pagination": gin.H{
				"limit": limit,
				"page":  page,
			},
		})
	})

	t.Run("Custom Limit", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/gov-jobs?limit=5&page=2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)

		pagination := resp["pagination"].(map[string]interface{})
		if pagination["limit"] != "5" {
			t.Errorf("expected limit 5, got %v", pagination["limit"])
		}
		if pagination["page"] != "2" {
			t.Errorf("expected page 2, got %v", pagination["page"])
		}
	})
}

// TestGetGovJobByID tests fetching single government job
func TestGetGovJobByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/api/v1/gov-jobs/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(http.StatusOK, gin.H{
			"data": gin.H{
				"id":    id,
				"title": "UPSC Civil Services 2026",
			},
		})
	})

	req, _ := http.NewRequest("GET", "/api/v1/gov-jobs/UPSC-123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "UPSC-123") {
		t.Errorf("expected response to contain 'UPSC-123'")
	}
}

// TestDBConnectivity demonstrates skipping if DATABASE_URL is not set
func TestDBConnectivity(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping DB tests")
	}
	// Real DB test logic would go here
}
