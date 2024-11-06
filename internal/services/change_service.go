package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type ChangeService struct {
	Repo *repositories.ChangeRepository
}

func (s *ChangeService) CreateChange(ctx context.Context, change *models.Change) (int, error) {
	return s.Repo.CreateChange(ctx, change)
}

func (s *ChangeService) GetChangeByID(ctx context.Context, id int) (*models.Change, error) {
	return s.Repo.GetChangeByID(ctx, id)
}

func (s *ChangeService) UpdateChange(ctx context.Context, change *models.Change) (*models.Change, error) {
	return s.Repo.UpdateChange(ctx, change)
}

func (s *ChangeService) DeleteChange(ctx context.Context, id int) error {
	return s.Repo.DeleteChange(ctx, id)
}

func (s *ChangeService) GetAllChangesByTransactionID(ctx context.Context, transactionID int) ([]models.Change, error) {
	return s.Repo.GetAllChangesByTransactionID(ctx, transactionID)
}
