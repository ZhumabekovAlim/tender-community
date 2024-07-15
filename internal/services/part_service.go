package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type PartService struct {
	repo *repositories.PartRepository
}

func NewPartService(repo *repositories.PartRepository) *PartService {
	return &PartService{repo: repo}
}

func (s *PartService) GetAllParts(ctx context.Context) ([]models.Part, error) {
	return s.repo.GetAllParts(ctx)
}

func (s *PartService) AddPart(ctx context.Context, part models.Part) error {
	// Add business logic validations if necessary
	return s.repo.AddPart(ctx, part)
}
