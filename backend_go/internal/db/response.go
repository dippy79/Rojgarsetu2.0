package db

// ErrorResponse creates a standardized error map for API responses.
func ErrorResponse(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"status": "error",
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

// SuccessResponse creates a standardized success map for API responses.
func SuccessResponse(data interface{}, pagination *Pagination) map[string]interface{} {
	resp := map[string]interface{}{
		"status": "success",
		"data":   data,
	}
	if pagination != nil {
		resp["pagination"] = pagination
	}
	return resp
}

// NewPagination creates a Pagination struct for API responses.
func NewPagination(page, limit, total int) Pagination {
	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}
}

// Pagination represents pagination metadata.
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
