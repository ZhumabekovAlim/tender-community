package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type TransactionHandler struct {
	Service *services.TransactionService
}

// CreateTransaction creates a new transaction.
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.CreateTransaction(r.Context(), transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetTransactionByID retrieves a transaction by ID.
func (h *TransactionHandler) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.Service.GetTransactionByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrTransactionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// GetAllTransactions retrieves all transactions.
func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.Service.GetAllTransactions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// UpdateTransaction updates an existing transaction.
func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err = json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	transaction.ID = id

	err = h.Service.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		if errors.Is(err, models.ErrTransactionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteTransaction deletes a transaction by ID.
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get(":id")
	if idStr == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteTransaction(r.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrTransactionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
