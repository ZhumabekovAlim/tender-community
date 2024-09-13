package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type BalanceHistoryHandler struct {
	Service *services.BalanceHistoryService
}

// CreateBalanceHistory creates a new balance history record.
func (h *BalanceHistoryHandler) CreateBalanceHistory(w http.ResponseWriter, r *http.Request) {
	var history models.BalanceHistory
	err := json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdHistory, err := h.Service.CreateBalanceHistory(r.Context(), history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdHistory)
}

// DeleteBalanceHistory deletes a balance history record by ID.
func (h *BalanceHistoryHandler) DeleteBalanceHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing balance history ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid balance history ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteBalanceHistory(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateBalanceHistory updates an existing balance history record.
func (h *BalanceHistoryHandler) UpdateBalanceHistory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing balance history ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid balance history ID", http.StatusBadRequest)
		return
	}

	var history models.BalanceHistory
	err = json.NewDecoder(r.Body).Decode(&history)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	history.ID = id

	updatedHistory, err := h.Service.UpdateBalanceHistory(r.Context(), history)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedHistory)
}

// GetAllBalanceHistories retrieves all balance history records by user_id.
func (h *BalanceHistoryHandler) GetBalanceHistoryByUserID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing balance history user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid balance history user ID", http.StatusBadRequest)
		return
	}

	histories, err := h.Service.GetBalanceHistoryByUserID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(histories)
}
