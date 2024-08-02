package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type TransactionService struct {
	Repo *repositories.TransactionRepository
}

// CreateTransaction creates a new transaction.
func (s *TransactionService) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	return s.Repo.CreateTransaction(ctx, transaction)
}

// GetTransactionByID retrieves a transaction by ID.
func (s *TransactionService) GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	return s.Repo.GetTransactionByID(ctx, id)
}

// GetAllTransactions retrieves all transactions.
func (s *TransactionService) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	return s.Repo.GetAllTransactions(ctx)
}

// UpdateTransaction updates an existing transaction.
func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {
	return s.Repo.UpdateTransaction(ctx, transaction)
}

// DeleteTransaction deletes a transaction by ID.
func (s *TransactionService) DeleteTransaction(ctx context.Context, id int) error {
	return s.Repo.DeleteTransaction(ctx, id)
}
