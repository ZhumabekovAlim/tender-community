package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type PersonalDebtService struct {
	Repo *repositories.PersonalDebtRepository
}

func (s *PersonalDebtService) CreatePersonalDebt(ctx context.Context, debt *models.PersonalDebt) (int, error) {
	return s.Repo.CreatePersonalDebt(ctx, debt)
}

func (s *PersonalDebtService) GetPersonalDebtByID(ctx context.Context, id int) (*models.PersonalDebt, error) {
	return s.Repo.GetPersonalDebtByID(ctx, id)
}

func (s *PersonalDebtService) UpdatePersonalDebt(ctx context.Context, debt *models.PersonalDebt) (*models.PersonalDebt, error) {
	return s.Repo.UpdatePersonalDebt(ctx, debt)
}

func (s *PersonalDebtService) DeletePersonalDebt(ctx context.Context, id int) error {
	return s.Repo.DeletePersonalDebt(ctx, id)
}

func (s *PersonalDebtService) GetAllPersonalDebts(ctx context.Context) ([]models.PersonalDebt, error) {
	return s.Repo.GetAllPersonalDebts(ctx)
}

func (s *PersonalDebtService) GetAllPersonalDebtsByStatus(ctx context.Context, id int) ([]models.PersonalDebt, error) {
	return s.Repo.GetAllPersonalDebtsByStatus(ctx, id)
}
