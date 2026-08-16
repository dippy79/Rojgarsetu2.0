package middleware

import (
	"html"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// SanitizeInput middleware sanitizes user input to prevent XSS attacks
func SanitizeInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				c.Request.URL.Query()[key][i] = sanitizeString(value)
			}
		}

		// Sanitize form data
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			c.Request.ParseForm()
			for key, values := range c.Request.PostForm {
				for i, value := range values {
					c.Request.PostForm[key][i] = sanitizeString(value)
				}
			}
		}

		c.Next()
	}
}

// sanitizeString removes potentially dangerous characters and HTML
func sanitizeString(input string) string {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Escape HTML entities to prevent XSS
	input = html.EscapeString(input)
	
	// Remove potential script tags (additional protection)
	input = strings.ReplaceAll(input, "<script", "")
	input = strings.ReplaceAll(input, "</script>", "")
	input = strings.ReplaceAll(input, "javascript:", "")
	input = strings.ReplaceAll(input, "onerror=", "")
	input = strings.ReplaceAll(input, "onload=", "")
	
	return input
}

// ValidateContentType ensures the request content type is valid
func ValidateContentType(contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			ct := c.GetHeader("Content-Type")
			if ct == "" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Content-Type header required"})
				c.Abort()
				return
			}
			
			// Check if content type matches expected type
			if !strings.HasPrefix(ct, contentType) {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Invalid Content-Type"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ValidateRequiredFields middleware checks for required fields in request body
func ValidateRequiredFields(fields []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			var body map[string]interface{}
			if err := c.ShouldBindJSON(&body); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
				c.Abort()
				return
			}
			
			for _, field := range fields {
				if _, exists := body[field]; !exists {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required field: " + field})
					c.Abort()
					return
				}
			}
		}
		c.Next()
	}
}