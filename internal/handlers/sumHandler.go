package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tender/internal/services"
)

type SumHandler struct {
	Service *services.SumService
}

// GetSumsByUserID handles the HTTP request to get sums by user_id.
func (h *SumHandler) GetSumsByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get(":id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id parameter", http.StatusBadRequest)
		return
	}

	sums, err := h.Service.GetSumsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sums)
}
