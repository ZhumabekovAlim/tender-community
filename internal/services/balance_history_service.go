package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type BalanceHistoryService struct {
	Repo *repositories.BalanceHistoryRepository
}

// CreateBalanceHistory creates a new balance history record.
func (s *BalanceHistoryService) CreateBalanceHistory(ctx context.Context, history models.BalanceHistory) (models.BalanceHistory, error) {
	id, err := s.Repo.CreateBalanceHistory(ctx, history)
	if err != nil {
		return models.BalanceHistory{}, err
	}

	history.ID = id
	return history, nil
}

// DeleteBalanceHistory deletes a balance history record by ID.
func (s *BalanceHistoryService) DeleteBalanceHistory(ctx context.Context, id int) error {
	return s.Repo.DeleteBalanceHistory(ctx, id)
}

// UpdateBalanceHistory updates an existing balance history record.
func (s *BalanceHistoryService) UpdateBalanceHistory(ctx context.Context, history models.BalanceHistory) (models.BalanceHistory, error) {
	return s.Repo.UpdateBalanceHistory(ctx, history)
}

// GetBalanceHistoryByID retrieves a balance history record by ID.
func (s *BalanceHistoryService) GetBalanceHistoryByUserID(ctx context.Context, id int) (models.BalanceHistory, error) {
	return s.Repo.GetBalanceHistoryByUserID(ctx, id)
}

// GetAllBalanceHistories retrieves all balance history records.
func (s *BalanceHistoryService) GetAllBalanceHistories(ctx context.Context) ([]models.BalanceHistory, error) {
	return s.Repo.GetAllBalanceHistories(ctx)
}
