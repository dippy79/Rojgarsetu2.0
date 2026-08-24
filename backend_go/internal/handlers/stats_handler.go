package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rojgarsetu/backend/internal/db"
)

type StatsHandler struct {
	db *db.PostgresDB
}

func NewStatsHandler(database *db.PostgresDB) *StatsHandler {
	return &StatsHandler{db: database}
}

func (h *StatsHandler) GetPlatformStats(c *gin.Context) {
	stats, err := h.db.Queries.GetPlatformStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
