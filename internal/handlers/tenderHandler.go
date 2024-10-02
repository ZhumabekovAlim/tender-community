package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type TenderHandler struct {
	Service *services.TenderService
}

// CreateTender creates a new tender.
func (h *TenderHandler) CreateTender(w http.ResponseWriter, r *http.Request) {
	var tender models.Tender
	err := json.NewDecoder(r.Body).Decode(&tender)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdTender, err := h.Service.CreateTender(r.Context(), tender)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTender)
}

// DeleteTender deletes a tender by ID.
func (h *TenderHandler) DeleteTender(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing tender ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteTender(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrTenderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateTender updates an existing tender.
func (h *TenderHandler) UpdateTender(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing tender ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	var tender models.Tender
	err = json.NewDecoder(r.Body).Decode(&tender)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tender.ID = id

	updatedTender, err := h.Service.UpdateTender(r.Context(), tender)
	if err != nil {
		if errors.Is(err, models.ErrTenderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTender)
}

// GetTenderByID retrieves a tender by ID.
func (h *TenderHandler) GetTenderByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing tender ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid tender ID", http.StatusBadRequest)
		return
	}

	tender, err := h.Service.GetTenderByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrTenderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tender)
}

// GetAllTenders retrieves all tenders.
func (h *TenderHandler) GetAllTenders(w http.ResponseWriter, r *http.Request) {
	tenders, err := h.Service.GetAllTenders(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tenders)
}
