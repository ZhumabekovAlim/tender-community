package handlers

import (
	"encoding/json"
	"net/http"
	"tender/internal/models"
	"tender/internal/services"
)

type PartsHandler struct {
	service *services.PartService
}

func NewPartsHandler(service *services.PartService) *PartsHandler {
	return &PartsHandler{service: service}
}

func (h *PartsHandler) GetAllParts(w http.ResponseWriter, r *http.Request) {
	parts, err := h.service.GetAllParts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parts)
}

func (h *PartsHandler) AddPart(w http.ResponseWriter, r *http.Request) {
	var part models.Part
	err := json.NewDecoder(r.Body).Decode(&part)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.AddPart(r.Context(), part); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
