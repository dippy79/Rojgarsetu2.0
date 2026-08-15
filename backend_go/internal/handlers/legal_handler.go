package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/legal"
)

type LegalHandler struct {
	db *sql.DB
}

func NewLegalHandler(db *sql.DB) *LegalHandler {
	return &LegalHandler{db: db}
}

// GET /api/v1/legal/disclaimer
func (h *LegalHandler) GetDisclaimer(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"disclaimer":       legal.GetDisclaimer(),
			"privacy_policy":   legal.GetPrivacyPolicy(),
			"terms_of_service": legal.GetTermsOfService(),
		},
	})
}

// POST /api/v1/legal/takedown
func (h *LegalHandler) PostTakedown(c *gin.Context) {
	var request struct {
		JobID     int    `json:"job_id"`
		FormID    int    `json:"form_id"`
		Requester string `json:"requester" binding:"required"`
		Reason    string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	query := `
		INSERT INTO takedown_requests (job_id, form_id, requester, reason, status)
		VALUES ($1, $2, $3, $4, 'PENDING')
	`

	_, err := h.db.Exec(query, request.JobID, request.FormID, request.Requester, request.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit takedown request: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Takedown request submitted successfully",
	})
}

// GET /api/v1/crawler/forms (proxy to crawler service)
func (h *LegalHandler) GetForms(c *gin.Context) {
	// Query forms from database
	rows, err := h.db.Query("SELECT id, title, conducting_body, form_type, official_website, is_taken_down FROM gov_forms_info ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch forms: " + err.Error()})
		return
	}
	defer rows.Close()

	var forms []gin.H
	for rows.Next() {
		var id int
		var title, conductingBody, formType, officialWebsite string
		var isTakenDown bool
		if err := rows.Scan(&id, &title, &conductingBody, &formType, &officialWebsite, &isTakenDown); err != nil {
			continue
		}
		forms = append(forms, gin.H{
			"id":               id,
			"title":            title,
			"conducting_body":  conductingBody,
			"form_type":        formType,
			"official_website": officialWebsite,
			"is_taken_down":    isTakenDown,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    forms,
		"count":   len(forms),
	})
}
