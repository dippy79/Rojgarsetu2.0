package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type VideoHandler struct {
	svc *services.ContentService
}

func NewVideoHandler(svc *services.ContentService) *VideoHandler {
	return &VideoHandler{svc: svc}
}

func (h *VideoHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *VideoHandler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}

func (h *VideoHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := db.VideoFilter{
		Channel:  r.URL.Query().Get("channel"),
		Category: r.URL.Query().Get("category"),
	}

	rows, count, err := h.svc.GetVideos(r.Context(), filter, page, limit)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]any{"data": rows, "count": count})
}

func (h *VideoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row, err := h.svc.GetVideoByID(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, row)
}
