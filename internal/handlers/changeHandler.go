package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type ChangeHandler struct {
	Service *services.ChangeService
}

func NewChangeHandler(service *services.ChangeService) *ChangeHandler {
	return &ChangeHandler{Service: service}
}

// Create a new change
func (h *ChangeHandler) CreateChange(w http.ResponseWriter, r *http.Request) {
	var change models.Change
	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateChange(r.Context(), &change)
	if err != nil {
		log.Printf("Error creating change: %v", err)
		http.Error(w, "Failed to create change", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Get a change by ID
func (h *ChangeHandler) GetChangeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing change ID", http.StatusBadRequest)
		return
	}

	change, err := h.Service.GetChangeByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching change: %v", err)
		http.Error(w, "Failed to fetch change", http.StatusInternalServerError)
		return
	}
	if change == nil {
		http.Error(w, "Change not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(change)
}

func (h *ChangeHandler) UpdateChange(w http.ResponseWriter, r *http.Request) {
	var change models.Change
	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedChange, err := h.Service.UpdateChange(r.Context(), &change)
	if err != nil {
		log.Printf("Error updating change: %v", err)
		http.Error(w, "Failed to update change", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedChange)
}

// Delete a change by ID
func (h *ChangeHandler) DeleteChange(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing change ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteChange(r.Context(), id); err != nil {
		log.Printf("Error deleting change: %v", err)
		http.Error(w, "Failed to delete change", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ChangeHandler) GetAllChangesByTransactionID(w http.ResponseWriter, r *http.Request) {
	transactionIDStr := r.URL.Query().Get(":transaction_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil || transactionIDStr == "" {
		http.Error(w, "Invalid or missing transaction_id", http.StatusBadRequest)
		return
	}

	changes, err := h.Service.GetAllChangesByTransactionID(r.Context(), transactionID)
	if err != nil {
		log.Printf("Error fetching changes by transaction_id: %v", err)
		http.Error(w, "Failed to fetch changes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(changes)
}
