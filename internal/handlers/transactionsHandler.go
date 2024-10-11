package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tender/internal/models"
	"tender/internal/services"
)

type TransactionHandler struct {
	Service                 *services.TransactionService
	ExtraTransactionService *services.ExtraTransactionService
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
	// Fetch regular transactions
	transactions, err := h.Service.GetAllTransactions(r.Context())
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	// Fetch extra transactions using the ExtraTransactionService
	extraTransactions, err := h.ExtraTransactionService.GetAllExtraTransactions(r.Context())
	if err != nil {
		log.Printf("Error fetching extra transactions: %v", err)
		http.Error(w, "Failed to fetch extra transactions", http.StatusInternalServerError)
		return
	}
	// Combine and send response
	combinedTransactions := combineTransactions(transactions, extraTransactions)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(combinedTransactions); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func combineTransactions(transactions []models.Transaction, extraTransactions []models.ExtraTransaction) []interface{} {
	var combined []interface{}

	// Add regular transactions to the combined slice
	for _, t := range transactions {
		combined = append(combined, t)
	}

	// Add extra transactions to the combined slice
	for _, et := range extraTransactions {
		combined = append(combined, et)
	}

	return combined
}

func (h *TransactionHandler) GetTransactionsByUser(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.Service.GetTransactionsByUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransactionsByCompany(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.Service.GetTransactionsByCompany(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransactionsForUserByCompany(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get(":user_id")

	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	companyIDStr := r.URL.Query().Get(":company_id")

	if companyIDStr == "" {
		http.Error(w, "Missing company ID", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.Service.GetTransactionsForUserByCompany(r.Context(), userID, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetAllTransactionsSum(w http.ResponseWriter, r *http.Request) {

	transactions, err := h.Service.GetAllTransactionsSum(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransactionCountsByUserID(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.Service.GetTransactionCountsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (h *TransactionHandler) GetTransactionsDebtZakup(w http.ResponseWriter, r *http.Request) {
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

	transactions, err := h.Service.GetTransactionsDebtZakup(r.Context(), userID)
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
func (h *TransactionHandler) GetMonthlyAmountsByGlobal(w http.ResponseWriter, r *http.Request) {
	monthlyAmounts, err := h.Service.GetMonthlyAmountsByGlobal(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 2
func (h *TransactionHandler) GetMonthlyAmountsByYear(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "Missing year", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsByYear(r.Context(), year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

func (h *TransactionHandler) GetMonthlyAmountsByCompany(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get("company_id")

	if companyIDStr == "" {
		http.Error(w, "Missing company ID", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsByCompany(r.Context(), companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 3
func (h *TransactionHandler) GetMonthlyAmountsByYearAndCompany(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	companyIDStr := r.URL.Query().Get("company_id")

	if yearStr == "" || companyIDStr == "" {
		http.Error(w, "Missing year or company ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsByYearAndCompany(r.Context(), year, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 4
func (h *TransactionHandler) GetMonthlyAmountsGroupedByYear(w http.ResponseWriter, r *http.Request) {
	monthlyAmounts, err := h.Service.GetMonthlyAmountsGroupedByYear(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 5
func (h *TransactionHandler) GetMonthlyAmountsGroupedByYearForUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsGroupedByYearForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 6
func (h *TransactionHandler) GetMonthlyAmountsForUserByYear(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")

	if userIDStr == "" || yearStr == "" {
		http.Error(w, "Missing user ID or year", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsForUserByYear(r.Context(), userID, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 7
func (h *TransactionHandler) GetMonthlyAmountsForUserByYearAndCompany(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")
	companyIDStr := r.URL.Query().Get("company_id")

	if userIDStr == "" || yearStr == "" || companyIDStr == "" {
		http.Error(w, "Missing user ID, year, or company ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	monthlyAmounts, err := h.Service.GetMonthlyAmountsForUserByYearAndCompany(r.Context(), userID, year, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(monthlyAmounts)
}

// 8
func (h *TransactionHandler) GetTotalAmountGroupedByCompany(w http.ResponseWriter, r *http.Request) {
	totalAmounts, err := h.Service.GetTotalAmountGroupedByCompany(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 9
func (h *TransactionHandler) GetTotalAmountByCompanyForYear(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	if yearStr == "" {
		http.Error(w, "Missing year", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForYear(r.Context(), year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

func (h *TransactionHandler) GetTotalAmountByCompanyForMonth(w http.ResponseWriter, r *http.Request) {
	monthStr := r.URL.Query().Get("month")

	if monthStr == "" {
		http.Error(w, "Missing month", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForMonth(r.Context(), month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 10
func (h *TransactionHandler) GetTotalAmountByCompanyForYearAndMonth(w http.ResponseWriter, r *http.Request) {
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	if yearStr == "" || monthStr == "" {
		http.Error(w, "Missing year or month", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForYearAndMonth(r.Context(), year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 11
func (h *TransactionHandler) GetTotalAmountGroupedByCompanyForUsers(w http.ResponseWriter, r *http.Request) {
	totalAmounts, err := h.Service.GetTotalAmountGroupedByCompanyForUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 12
func (h *TransactionHandler) GetTotalAmountByCompanyForUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

func (h *TransactionHandler) GetTotalAmountByCompanyForUserAndMonth(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	monthStr := r.URL.Query().Get("month")

	if userIDStr == "" || monthStr == "" {
		http.Error(w, "Missing user ID or month", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForUserAndMonth(r.Context(), userID, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 13
func (h *TransactionHandler) GetTotalAmountByCompanyForUserAndYear(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")

	if userIDStr == "" || yearStr == "" {
		http.Error(w, "Missing user ID or year", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForUserAndYear(r.Context(), userID, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

// 14
func (h *TransactionHandler) GetTotalAmountByCompanyForUserYearAndMonth(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

	if userIDStr == "" || yearStr == "" || monthStr == "" {
		http.Error(w, "Missing user ID, year, or month", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		http.Error(w, "Invalid year", http.StatusBadRequest)
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	totalAmounts, err := h.Service.GetTotalAmountByCompanyForUserYearAndMonth(r.Context(), userID, year, month)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(totalAmounts)
}

func (h *TransactionHandler) GetAllByUserIDAndStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user_id and status from the query parameters
	userIDStr := r.URL.Query().Get(":user_id")
	statusStr := r.URL.Query().Get(":status")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userIDStr == "" {
		http.Error(w, "Invalid or missing user_id", http.StatusBadRequest)
		return
	}

	status, err := strconv.Atoi(statusStr)
	if err != nil || statusStr == "" {
		http.Error(w, "Invalid or missing status", http.StatusBadRequest)
		return
	}

	// Call the service to get all data by user ID and status
	allData, err := h.Service.GetAllByUserIDAndStatus(r.Context(), userID, status)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}

	// Write the combined result as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allData)
}

func (h *TransactionHandler) GetAllTransactionsByDateRange(w http.ResponseWriter, r *http.Request) {
	var dateRange models.DateRangeRequest

	// Parse the request body to get the date range
	if err := json.NewDecoder(r.Body).Decode(&dateRange); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch transactions within the date range
	transactions, err := h.Service.GetAllTransactionsByDateRange(r.Context(), dateRange.StartDate, dateRange.EndDate, dateRange.UserId)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TransactionHandler) GetAllTransactionsByDateRangeCompany(w http.ResponseWriter, r *http.Request) {
	var dateRange models.DateRangeRequestCompany

	// Parse the request body to get the date range
	if err := json.NewDecoder(r.Body).Decode(&dateRange); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch transactions within the date range
	transactions, err := h.Service.GetAllTransactionsByDateRangeCompany(r.Context(), dateRange.StartDate, dateRange.EndDate, dateRange.UserId, dateRange.CompanyId)
	if err != nil {
		log.Printf("Error fetching transactions: %v", err)
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}

	// Send the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
