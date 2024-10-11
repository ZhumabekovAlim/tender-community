package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type TransactionService struct {
	Repo                 *repositories.TransactionRepository
	ExtraTransactionRepo *repositories.ExtraTransactionRepository
}

type CombinedTransactions struct {
	Transactions      []models.Transaction      `json:"transactions"`
	Tenders           []models.Tender           `json:"tenders"`
	ExtraTransactions []models.ExtraTransaction `json:"extra_transactions"`
}

// CreateTransaction creates a new transaction with expenses.
func (s *TransactionService) CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return s.Repo.CreateTransaction(ctx, transaction)
}

// GetTransactionByID retrieves a transaction by ID along with its expenses.
func (s *TransactionService) GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	return s.Repo.GetTransactionByID(ctx, id)
}

// GetAllTransactions retrieves all transactions.
func (s *TransactionService) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	return s.Repo.GetAllTransactions(ctx)
}

// GetTransactionsByUser retrieves transactions by user id.
func (s *TransactionService) GetTransactionsByUser(ctx context.Context, userID int) ([]models.Transaction, error) {
	return s.Repo.GetTransactionsByUser(ctx, userID)
}

// GetTransactionsByCompany retrieves transactions by company id.
func (s *TransactionService) GetTransactionsByCompany(ctx context.Context, companyID int) ([]models.Transaction, error) {
	return s.Repo.GetTransactionsByCompany(ctx, companyID)
}

// GetTransactionsByCompany retrieves transactions by company id.
func (s *TransactionService) GetTransactionsForUserByCompany(ctx context.Context, userID, companyID int) ([]models.Transaction, error) {
	return s.Repo.GetTransactionsForUserByCompany(ctx, userID, companyID)
}

func (s *TransactionService) GetAllTransactionsSum(ctx context.Context) (*models.TransactionDebt, error) {
	return s.Repo.GetAllTransactionsSum(ctx)
}

func (s *TransactionService) GetTransactionCountsByUserID(ctx context.Context, userID int) (*models.TransactionCount, error) {
	// Call the repository function to get transaction counts by user ID
	return s.Repo.GetTransactionCountsByUserID(ctx, userID)
}

func (s *TransactionService) GetTransactionsDebtZakup(ctx context.Context, userID int) (*models.TransactionDebt, error) {
	return s.Repo.GetTransactionsDebtZakup(ctx, userID)
}

// UpdateTransaction updates an existing transaction and its expenses.
func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return s.Repo.UpdateTransaction(ctx, transaction)
}

// DeleteTransaction deletes a transaction and its expenses by ID.
func (s *TransactionService) DeleteTransaction(ctx context.Context, id int) error {
	return s.Repo.DeleteTransaction(ctx, id)
}

func (s *TransactionService) GetMonthlyAmountsByGlobal(ctx context.Context) ([]repositories.YearlyAmounts, error) {
	return s.Repo.GetMonthlyAmountsByGlobal(ctx)
}

func (s *TransactionService) GetMonthlyAmountsByYear(ctx context.Context, year int) ([]repositories.MonthlyAmount, error) {
	return s.Repo.GetMonthlyAmountsByYear(ctx, year)
}

func (s *TransactionService) GetMonthlyAmountsByCompany(ctx context.Context, companyID int) ([]repositories.MonthlyAmount, error) {
	return s.Repo.GetMonthlyAmountsByCompany(ctx, companyID)
}

func (s *TransactionService) GetMonthlyAmountsByYearAndCompany(ctx context.Context, year int, companyID int) ([]repositories.MonthlyAmount, error) {
	return s.Repo.GetMonthlyAmountsByYearAndCompany(ctx, year, companyID)
}

func (s *TransactionService) GetMonthlyAmountsGroupedByYear(ctx context.Context) ([]repositories.YearlyAmounts, error) {
	return s.Repo.GetMonthlyAmountsGroupedByYear(ctx)
}

func (s *TransactionService) GetMonthlyAmountsGroupedByYearForUser(ctx context.Context, userID int) ([]repositories.YearlyAmounts, error) {
	return s.Repo.GetMonthlyAmountsGroupedByYearForUser(ctx, userID)
}

func (s *TransactionService) GetMonthlyAmountsForUserByYear(ctx context.Context, userID int, year int) ([]repositories.MonthlyAmount, error) {
	return s.Repo.GetMonthlyAmountsForUserByYear(ctx, userID, year)
}

func (s *TransactionService) GetMonthlyAmountsForUserByYearAndCompany(ctx context.Context, userID int, year int, companyID int) ([]repositories.MonthlyAmount, error) {
	return s.Repo.GetMonthlyAmountsForUserByYearAndCompany(ctx, userID, year, companyID)
}

func (s *TransactionService) GetTotalAmountGroupedByCompany(ctx context.Context) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountGroupedByCompany(ctx)
}

func (s *TransactionService) GetTotalAmountByCompanyForYear(ctx context.Context, year int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForYear(ctx, year)
}

func (s *TransactionService) GetTotalAmountByCompanyForMonth(ctx context.Context, month int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForMonth(ctx, month)
}

func (s *TransactionService) GetTotalAmountByCompanyForYearAndMonth(ctx context.Context, year int, month int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForYearAndMonth(ctx, year, month)
}

func (s *TransactionService) GetTotalAmountGroupedByCompanyForUsers(ctx context.Context) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountGroupedByCompany(ctx)
}

func (s *TransactionService) GetTotalAmountByCompanyForUser(ctx context.Context, userID int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForUser(ctx, userID)
}

func (s *TransactionService) GetTotalAmountByCompanyForUserAndMonth(ctx context.Context, userID int, month int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForUserAndMonth(ctx, userID, month)
}

func (s *TransactionService) GetTotalAmountByCompanyForUserAndYear(ctx context.Context, userID int, year int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForUserAndYear(ctx, userID, year)
}

func (s *TransactionService) GetTotalAmountByCompanyForUserYearAndMonth(ctx context.Context, userID int, year int, month int) ([]repositories.CompanyTotalAmount, error) {
	return s.Repo.GetTotalAmountByCompanyForUserYearAndMonth(ctx, userID, year, month)
}

func (s *TransactionService) GetAllByUserIDAndStatus(ctx context.Context, userID, status int) (*CombinedTransactions, error) {
	// Fetch transactions, tenders, and extra_transactions
	transactions, err := s.Repo.FindAllTransactionsByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	tenders, err := s.Repo.FindAllTendersByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	extraTransactions, err := s.Repo.FindAllExtraTransactionsByUserIDAndStatus(ctx, userID, status)
	if err != nil {
		return nil, err
	}

	// Return all data in a single struct
	return &CombinedTransactions{
		Transactions:      transactions,
		Tenders:           tenders,
		ExtraTransactions: extraTransactions,
	}, nil
}

func (s *TransactionService) GetAllTransactionsByDateRange(ctx context.Context, startDate, endDate string) ([]models.Transaction, error) {
	return s.Repo.GetAllTransactionsByDateRange(ctx, startDate, endDate)
}
