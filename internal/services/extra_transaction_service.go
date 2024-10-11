package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type ExtraTransactionService struct {
	Repo *repositories.ExtraTransactionRepository
}

func (s *ExtraTransactionService) CreateExtraTransaction(ctx context.Context, extraTransaction models.ExtraTransaction) (models.ExtraTransaction, error) {
	return s.Repo.CreateExtraTransaction(ctx, extraTransaction)
}

func (s *ExtraTransactionService) GetExtraTransactionByID(ctx context.Context, id int) (models.ExtraTransaction, error) {
	return s.Repo.GetExtraTransactionByID(ctx, id)
}

func (s *ExtraTransactionService) GetAllExtraTransactions(ctx context.Context) ([]models.ExtraTransaction, error) {
	return s.Repo.GetAllExtraTransactions(ctx)
}

func (s *ExtraTransactionService) GetExtraTransactionsByUser(ctx context.Context, userID int) ([]models.ExtraTransaction, error) {
	return s.Repo.GetExtraTransactionsByUser(ctx, userID)
}

func (s *ExtraTransactionService) UpdateExtraTransaction(ctx context.Context, extraTransaction models.ExtraTransaction) (models.ExtraTransaction, error) {
	return s.Repo.UpdateExtraTransaction(ctx, extraTransaction)
}

func (s *ExtraTransactionService) DeleteExtraTransaction(ctx context.Context, id int) error {
	return s.Repo.DeleteExtraTransaction(ctx, id)
}

func (s *ExtraTransactionService) GetExtraTransactionCountsByUserID(ctx context.Context, userID int) (*models.ExtraTransactionCount, error) {
	return s.Repo.GetExtraTransactionCountsByUserID(ctx, userID)
}
