package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type TenderService struct {
	Repo *repositories.TenderRepository
}

// CreateTender creates a new tender.
func (s *TenderService) CreateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	id, err := s.Repo.CreateTender(ctx, tender)
	if err != nil {
		return models.Tender{}, err
	}

	tender.ID = id
	return tender, nil
}

// DeleteTender deletes a tender by ID.
func (s *TenderService) DeleteTender(ctx context.Context, id int) error {
	return s.Repo.DeleteTender(ctx, id)
}

// UpdateTender updates an existing tender.
func (s *TenderService) UpdateTender(ctx context.Context, tender models.Tender) (models.Tender, error) {
	return s.Repo.UpdateTender(ctx, tender)
}

// GetTenderByID retrieves a tender by ID.
func (s *TenderService) GetTenderByID(ctx context.Context, id int) (models.Tender, error) {
	return s.Repo.GetTenderByID(ctx, id)
}

// GetAllTenders retrieves all tenders.
func (s *TenderService) GetAllTenders(ctx context.Context) ([]models.Tender, error) {
	return s.Repo.GetAllTenders(ctx)
}

func (s *TenderService) GetTendersByUserID(ctx context.Context, userID int) ([]models.Tender, error) {
	// Add any business logic here (e.g., validation) if needed before calling the repository.
	return s.Repo.GetTendersByUserID(ctx, userID)
}
