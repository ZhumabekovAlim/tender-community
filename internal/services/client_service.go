package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type ClientService struct {
	Repo *repositories.ClientRepository
}

// GetClientData retrieves the data for a specific client
func (s *ClientService) GetClientData(ctx context.Context, userID int) (models.ClientData, error) {
	return s.Repo.GetClientData(ctx, userID)
}
