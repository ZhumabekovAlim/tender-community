package services

import (
	"context"
	"tender/internal/models"
)

// AddPermission adds a new permission for a user.
func (s *UserService) AddPermission(ctx context.Context, permission models.Permission) error {
	return s.Repo.AddPermission(ctx, permission)
}

// DeletePermission deletes a permission for a user.
func (s *UserService) DeletePermission(ctx context.Context, id int) error {
	return s.Repo.DeletePermission(ctx, id)
}

// UpdatePermission updates an existing permission.
func (s *UserService) UpdatePermission(ctx context.Context, permission models.Permission) error {
	return s.Repo.UpdatePermission(ctx, permission)
}

// GetPermission retrieves a permission by ID.
func (s *UserService) GetPermissionsByUserID(ctx context.Context, userID int) ([]models.Permission, error) {
	return s.Repo.GetPermissionsByUserID(ctx, userID)
}
