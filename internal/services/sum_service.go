package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type SumService struct {
	Repo *repositories.SumRepository
}

// GetSumsByUserID calls the repository to get the sums.
func (s *SumService) GetSumsByUserID(ctx context.Context, userID int) (models.Sums, error) {
	return s.Repo.GetSumsByUserID(ctx, userID)
}
