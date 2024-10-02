package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tender/internal/services"
)

type ClientHandler struct {
	Service *services.ClientService
}

// GetClientData handles the HTTP request to get data for a specific client
func (h *ClientHandler) GetClientData(w http.ResponseWriter, r *http.Request) {
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

	clientData, err := h.Service.GetClientData(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clientData)
}
