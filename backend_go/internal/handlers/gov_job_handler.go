package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rojgarsetu/backend/internal/db"
	"github.com/rojgarsetu/backend/internal/services"
)

type GovJobHandler struct {
	svc *services.ContentService
}

func NewGovJobHandler(svc *services.ContentService) *GovJobHandler {
	return &GovJobHandler{svc: svc}
}

func (h *GovJobHandler) writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *GovJobHandler) writeError(w http.ResponseWriter, status int, msg string) {
	h.writeJSON(w, status, map[string]string{"error": msg})
}

func (h *GovJobHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := db.GovJobFilter{
		Department: r.URL.Query().Get("department"),
		Location:   r.URL.Query().Get("location"),
		Source:     r.URL.Query().Get("source"),
	}

	rows, count, err := h.svc.GetGovJobs(r.Context(), filter, page, limit)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]any{"data": rows, "count": count})
}

func (h *GovJobHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	row, err := h.svc.GetGovJobByID(r.Context(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, row)
}
