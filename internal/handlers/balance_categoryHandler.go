package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type BalanceCategoryHandler struct {
	Service *services.BalanceCategoryService
}

func NewBalanceCategoryHandler(service *services.BalanceCategoryService) *BalanceCategoryHandler {
	return &BalanceCategoryHandler{Service: service}
}

// Create a new balance category
func (h *BalanceCategoryHandler) CreateBalanceCategory(w http.ResponseWriter, r *http.Request) {
	var category models.BalanceCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.Service.CreateBalanceCategory(r.Context(), &category)
	if err != nil {
		log.Printf("Error creating balance category: %v", err)
		http.Error(w, "Failed to create balance category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

// Get a balance category by ID
func (h *BalanceCategoryHandler) GetBalanceCategoryByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing balance category ID", http.StatusBadRequest)
		return
	}

	category, err := h.Service.GetBalanceCategoryByID(r.Context(), id)
	if err != nil {
		log.Printf("Error fetching balance category: %v", err)
		http.Error(w, "Failed to fetch balance category", http.StatusInternalServerError)
		return
	}
	if category == nil {
		http.Error(w, "Balance category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// Update a balance category by ID
func (h *BalanceCategoryHandler) UpdateBalanceCategory(w http.ResponseWriter, r *http.Request) {
	var category models.BalanceCategory
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedCategory, err := h.Service.UpdateBalanceCategory(r.Context(), &category)
	if err != nil {
		log.Printf("Error updating balance category: %v", err)
		http.Error(w, "Failed to update balance category", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCategory)
}

// Delete a balance category by ID
func (h *BalanceCategoryHandler) DeleteBalanceCategory(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil || idStr == "" {
		http.Error(w, "Invalid or missing balance category ID", http.StatusBadRequest)
		return
	}

	if err := h.Service.DeleteBalanceCategory(r.Context(), id); err != nil {
		log.Printf("Error deleting balance category: %v", err)
		http.Error(w, "Failed to delete balance category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get all balance categories
func (h *BalanceCategoryHandler) GetAllBalanceCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Service.GetAllBalanceCategories(r.Context())
	if err != nil {
		log.Printf("Error fetching all balance categories: %v", err)
		http.Error(w, "Failed to fetch balance categories", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
