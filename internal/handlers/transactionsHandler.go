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

// CreateTransaction creates a new transaction with expenses.
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdTransaction, err := h.Service.CreateTransaction(r.Context(), transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTransaction)
}

// GetTransactionByID retrieves a transaction by ID along with its expenses.
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

// UpdateTransaction updates an existing transaction and its expenses.
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

	updatedTransaction, err := h.Service.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		if errors.Is(err, models.ErrTransactionNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTransaction)
}

// DeleteTransaction deletes a transaction and its expenses by ID.
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

// 1
func (h *TransactionHandler) GetMonthlyAmountsByYear(w http.ResponseWriter, r *http.Request) {
	monthlyAmounts, err := h.Service.GetMonthlyAmountsByYear(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}
