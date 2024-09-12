package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type CategoryRepository struct {
	Db *sql.DB
}

// CreateCategory inserts a new category into the database.
func (r *CategoryRepository) CreateCategory(ctx context.Context, category models.Category) (int, error) {
	result, err := r.Db.ExecContext(ctx, "INSERT INTO categories (category_name) VALUES (?)", category.CategoryName)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// DeleteCategory removes a category from the database by ID.
func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrCategoryNotFound
	}

	return nil
}

// UpdateCategory updates an existing category in the database.
func (r *CategoryRepository) UpdateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	_, err := r.Db.ExecContext(ctx, "UPDATE categories SET category_name = ? WHERE id = ?", category.CategoryName, category.ID)
	if err != nil {
		return models.Category{}, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category by ID from the database.
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (models.Category, error) {
	var category models.Category
	err := r.Db.QueryRowContext(ctx, "SELECT id, category_name FROM categories WHERE id = ?", id).Scan(&category.ID, &category.CategoryName)
	if err != nil {
		if err == sql.ErrNoRows {
			return category, models.ErrCategoryNotFound
		}
		return category, err
	}

	return category, nil
}

// GetAllCategories retrieves all categories from the database.
func (r *CategoryRepository) GetAllCategories(ctx context.Context) ([]models.Category, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, category_name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.CategoryName)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
