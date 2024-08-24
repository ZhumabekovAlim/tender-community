package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type ExtraTransactionHandler struct {
	Service *services.ExtraTransactionService
}

func (h *ExtraTransactionHandler) CreateExtraTransaction(w http.ResponseWriter, r *http.Request) {
	var extraTransaction models.ExtraTransaction
	err := json.NewDecoder(r.Body).Decode(&extraTransaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdExtraTransaction, err := h.Service.CreateExtraTransaction(r.Context(), extraTransaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdExtraTransaction)
}

func (h *ExtraTransactionHandler) GetExtraTransactionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing extra transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid extra transaction ID", http.StatusBadRequest)
		return
	}

	extraTransaction, err := h.Service.GetExtraTransactionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(extraTransaction)
}

func (h *ExtraTransactionHandler) GetAllExtraTransactions(w http.ResponseWriter, r *http.Request) {
	extraTransactions, err := h.Service.GetAllExtraTransactions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(extraTransactions)
}

func (h *ExtraTransactionHandler) GetExtraTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get(":id")

	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	extraTransactions, err := h.Service.GetExtraTransactionsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(extraTransactions)
}

func (h *ExtraTransactionHandler) UpdateExtraTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing extra transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid extra transaction ID", http.StatusBadRequest)
		return
	}

	var extraTransaction models.ExtraTransaction
	err = json.NewDecoder(r.Body).Decode(&extraTransaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	extraTransaction.ID = id

	updatedExtraTransaction, err := h.Service.UpdateExtraTransaction(r.Context(), extraTransaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedExtraTransaction)
}

func (h *ExtraTransactionHandler) DeleteExtraTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing extra transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid extra transaction ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteExtraTransaction(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
