package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/services"
)

type HistoryHandler struct {
	Service *services.HistoryService
}

// GetAllHistory handles the request to retrieve all history.
func (h *HistoryHandler) GetAllHistory(w http.ResponseWriter, r *http.Request) {
	history, err := h.Service.GetAllHistory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
