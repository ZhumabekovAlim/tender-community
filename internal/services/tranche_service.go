package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type TrancheService struct {
	Repo *repositories.TrancheRepository
}

func (s *TrancheService) CreateTranche(ctx context.Context, tranche *models.Tranche) (int, error) {
	return s.Repo.CreateTranche(ctx, tranche)
}

func (s *TrancheService) GetTrancheByID(ctx context.Context, id int) (*models.Tranche, error) {
	return s.Repo.GetTrancheByID(ctx, id)
}

func (s *TrancheService) UpdateTranche(ctx context.Context, tranche *models.Tranche) (*models.Tranche, error) {
	return s.Repo.UpdateTranche(ctx, tranche)
}

func (s *TrancheService) DeleteTranche(ctx context.Context, id int) error {
	return s.Repo.DeleteTranche(ctx, id)
}

func (s *TrancheService) GetAllTranchesByTransactionID(ctx context.Context, transactionID int) ([]models.Tranche, error) {
	return s.Repo.GetAllTranchesByTransactionID(ctx, transactionID)
}
