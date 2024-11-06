package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type PersonalDebtHandler struct {
	Service *services.PersonalDebtService
}

func NewPersonalDebtHandler(service *services.PersonalDebtService) *PersonalDebtHandler {
	return &PersonalDebtHandler{Service: service}
}

// Create a new personal debt
func (h *PersonalDebtHandler) CreatePersonalDebt(w http.ResponseWriter, r *http.Request) {
	var debt models.PersonalDebt
	if err := json.NewDecoder(r.Body).Decode(&debt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreatePersonalDebt(r.Context(), &debt)
	if err != nil {
		log.Printf("Error creating personal debt: %v", err)
		http.Error(w, "Failed to create personal debt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Get a personal debt by ID
func (h *PersonalDebtHandler) GetPersonalDebtByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing personal debt ID", http.StatusBadRequest)
		return
	}

	debt, err := h.Service.GetPersonalDebtByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching personal debt: %v", err)
		http.Error(w, "Failed to fetch personal debt", http.StatusInternalServerError)
		return
	}
	if debt == nil {
		http.Error(w, "Personal debt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debt)
}

// Update a personal debt by ID
func (h *PersonalDebtHandler) UpdatePersonalDebt(w http.ResponseWriter, r *http.Request) {
	var debt models.PersonalDebt
	if err := json.NewDecoder(r.Body).Decode(&debt); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedDebt, err := h.Service.UpdatePersonalDebt(r.Context(), &debt)
	if err != nil {
		log.Printf("Error updating personal debt: %v", err)
		http.Error(w, "Failed to update personal debt", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDebt)
}

// Delete a personal debt by ID
func (h *PersonalDebtHandler) DeletePersonalDebt(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing personal debt ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeletePersonalDebt(r.Context(), id); err != nil {
		log.Printf("Error deleting personal debt: %v", err)
		http.Error(w, "Failed to delete personal debt", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get all personal debts
func (h *PersonalDebtHandler) GetAllPersonalDebts(w http.ResponseWriter, r *http.Request) {
	debts, err := h.Service.GetAllPersonalDebts(r.Context())
	if err != nil {
		log.Printf("Error fetching all personal debts: %v", err)
		http.Error(w, "Failed to fetch personal debts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(debts)
}
