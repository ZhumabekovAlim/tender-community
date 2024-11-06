package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type BalanceCategoryService struct {
	Repo *repositories.BalanceCategoryRepository
}

func (s *BalanceCategoryService) CreateBalanceCategory(ctx context.Context, category *models.BalanceCategory) (int, error) {
	return s.Repo.CreateBalanceCategory(ctx, category)
}

func (s *BalanceCategoryService) GetBalanceCategoryByID(ctx context.Context, id int) (*models.BalanceCategory, error) {
	return s.Repo.GetBalanceCategoryByID(ctx, id)
}

func (s *BalanceCategoryService) UpdateBalanceCategory(ctx context.Context, category *models.BalanceCategory) (*models.BalanceCategory, error) {
	return s.Repo.UpdateBalanceCategory(ctx, category)
}

func (s *BalanceCategoryService) DeleteBalanceCategory(ctx context.Context, id int) error {
	return s.Repo.DeleteBalanceCategory(ctx, id)
}

func (s *BalanceCategoryService) GetAllBalanceCategories(ctx context.Context) ([]models.BalanceCategory, error) {
	return s.Repo.GetAllBalanceCategories(ctx)
}
