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
func (r *CategoryRepository) CreateCategory(ctx context.Context, category models.Category) (models.Category, error) {
	// Insert the new category into the database
	result, err := r.Db.ExecContext(ctx, "INSERT INTO categories (category_name,  parent_id) VALUES (?, ?)", category.CategoryName, category.ParentID)
	if err != nil {
		return models.Category{}, err
	}

	// Get the ID of the inserted record
	id, err := result.LastInsertId()
	if err != nil {
		return models.Category{}, err
	}

	// Query the database to get the full category model
	var createdCategory models.Category
	err = r.Db.QueryRowContext(ctx, "SELECT id, category_name FROM categories WHERE id = ?", id).
		Scan(&createdCategory.ID, &createdCategory.CategoryName)
	if err != nil {
		return models.Category{}, err
	}

	// Return the full category model
	return createdCategory, nil
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
	rows, err := r.Db.QueryContext(ctx, "SELECT id, category_name FROM categories WHERE parent_id = 0")
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

func (r *CategoryRepository) GetAllCategoriesByParent(ctx context.Context, parentID int) ([]models.Category, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, category_name FROM categories WHERE parent_id = ?", parentID)
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
