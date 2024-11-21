package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"
)

type HistoryHandler struct {
	Service *services.HistoryService
}

func (h *HistoryHandler) GetAllHistory(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body
	var req models.HistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Limit == 0 {
		req.Limit = 50
	}
	if req.Offset == 0 {
		req.Offset = 0
	} else {
		req.Offset = req.Offset - 1
	}

	history, err := h.Service.GetAllHistory(r.Context(), req.Source, req.StartDate, req.EndDate, req.Limit, req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}
