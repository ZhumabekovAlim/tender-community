package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type HistoryService struct {
	Repo *repositories.HistoryRepository
}

// GetAllHistory retrieves combined history data from the repository.
func (s *HistoryService) GetAllHistory(ctx context.Context) ([]models.CombinedAction, error) {
	return s.Repo.GetAllHistory(ctx)
}
