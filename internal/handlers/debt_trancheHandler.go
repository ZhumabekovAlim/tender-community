package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type DebtTrancheHandler struct {
	Service *services.DebtTrancheService
}

func NewDebtTrancheHandler(service *services.DebtTrancheService) *DebtTrancheHandler {
	return &DebtTrancheHandler{Service: service}
}

// Create a new tranche
func (h *DebtTrancheHandler) CreateDebtTranche(w http.ResponseWriter, r *http.Request) {
	var tranche models.DebtTranche
	if err := json.NewDecoder(r.Body).Decode(&tranche); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateDebtTranche(r.Context(), &tranche)
	if err != nil {
		log.Printf("Error creating tranche: %v", err)
		http.Error(w, "Failed to create tranche", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Get a tranche by ID
func (h *DebtTrancheHandler) GetDebtTrancheByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing tranche ID", http.StatusBadRequest)
		return
	}

	tranche, err := h.Service.GetDebtTrancheByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching tranche: %v", err)
		http.Error(w, "Failed to fetch tranche", http.StatusInternalServerError)
		return
	}
	if tranche == nil {
		http.Error(w, "Tranche not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tranche)
}

func (h *DebtTrancheHandler) UpdateDebtTranche(w http.ResponseWriter, r *http.Request) {
	var tranche models.DebtTranche
	if err := json.NewDecoder(r.Body).Decode(&tranche); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedTranche, err := h.Service.UpdateDebtTranche(r.Context(), &tranche)
	if err != nil {
		log.Printf("Error updating tranche: %v", err)
		http.Error(w, "Failed to update tranche", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTranche)
}

// Delete a tranche by ID
func (h *DebtTrancheHandler) DeleteDebtTranche(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing tranche ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteDebtTranche(r.Context(), id); err != nil {
		log.Printf("Error deleting tranche: %v", err)
		http.Error(w, "Failed to delete tranche", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DebtTrancheHandler) GetAllDebtTranchesByTransactionID(w http.ResponseWriter, r *http.Request) {
	// Extract transaction_id from query parameters
	transactionIDStr := r.URL.Query().Get(":debt_id")
	transactionID, err := strconv.Atoi(transactionIDStr)
	if err != nil || transactionIDStr == "" {
		http.Error(w, "Invalid or missing transaction_id", http.StatusBadRequest)
		return
	}

	// Fetch all tranches for the given transaction_id
	tranches, err := h.Service.GetAllDebtTranchesByTransactionID(r.Context(), transactionID)
	if err != nil {
		log.Printf("Error fetching tranches by transaction_id: %v", err)
		http.Error(w, "Failed to fetch tranches", http.StatusInternalServerError)
		return
	}

	// Respond with the list of tranches as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tranches)
}
