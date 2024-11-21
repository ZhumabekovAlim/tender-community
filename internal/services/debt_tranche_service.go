package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type DebtTrancheService struct {
	Repo *repositories.DebtTrancheRepository
}

func (s *DebtTrancheService) CreateDebtTranche(ctx context.Context, tranche *models.DebtTranche) (int, error) {
	return s.Repo.CreateDebtTranche(ctx, tranche)
}

func (s *DebtTrancheService) GetDebtTrancheByID(ctx context.Context, id int) (*models.DebtTranche, error) {
	return s.Repo.GetDebtTrancheByID(ctx, id)
}

func (s *DebtTrancheService) UpdateDebtTranche(ctx context.Context, tranche *models.DebtTranche) (*models.DebtTranche, error) {
	return s.Repo.UpdateDebtTranche(ctx, tranche)
}

func (s *DebtTrancheService) DeleteDebtTranche(ctx context.Context, id int) error {
	return s.Repo.DeleteDebtTranche(ctx, id)
}

func (s *DebtTrancheService) GetAllDebtTranchesByTransactionID(ctx context.Context, transactionID int) ([]models.DebtTranche, error) {
	return s.Repo.GetAllDebtTranchesByDebtID(ctx, transactionID)
}
