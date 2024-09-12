package services

import (
	"context"
	"tender/internal/models"
	"tender/internal/repositories"
)

type CategoryService struct {
	Repo *repositories.CategoryRepository
}

// CreateCategory creates a new category.
func (s *CategoryService) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	category, _ = s.Repo.CreateCategory(ctx, category)

	return category, nil
}

// DeleteCategory deletes a category by ID.
func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.Repo.DeleteCategory(ctx, id)
}

// UpdateCategory updates an existing category.
func (s *CategoryService) UpdateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	return s.Repo.UpdateCategory(ctx, category)
}

// GetCategoryByID retrieves a category by ID.
func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (models.Category, error) {
	return s.Repo.GetCategoryByID(ctx, id)
}

// GetAllCategories retrieves all categories.
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	return s.Repo.GetAllCategories(ctx)
}
