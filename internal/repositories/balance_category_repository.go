package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
)

type BalanceCategoryRepository struct {
	Db *sql.DB
}

func (r *BalanceCategoryRepository) CreateBalanceCategory(ctx context.Context, category *models.BalanceCategory) (int, error) {
	query := `
		INSERT INTO balance_category (name, parent_id) 
		VALUES (?, ?)
	`
	result, err := r.Db.ExecContext(ctx, query, category.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to create balance category: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	return int(id), nil
}

func (r *BalanceCategoryRepository) GetBalanceCategoryByID(ctx context.Context, id int) (*models.BalanceCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM balance_category
		WHERE id = ?
	`
	var category models.BalanceCategory
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get balance category: %w", err)
	}
	return &category, nil
}

func (r *BalanceCategoryRepository) UpdateBalanceCategory(ctx context.Context, category *models.BalanceCategory) (*models.BalanceCategory, error) {
	query := `
		UPDATE balance_category
		SET name = ?
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, category.Name, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance category: %w", err)
	}

	// Retrieve the updated balance category
	return r.GetBalanceCategoryByID(ctx, category.ID)
}

func (r *BalanceCategoryRepository) DeleteBalanceCategory(ctx context.Context, id int) error {
	query := `
		DELETE FROM balance_category
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete balance category: %w", err)
	}
	return nil
}

func (r *BalanceCategoryRepository) GetAllBalanceCategories(ctx context.Context) ([]models.BalanceCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM balance_category
		ORDER BY name
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query balance categories: %w", err)
	}
	defer rows.Close()

	var categories []models.BalanceCategory
	for rows.Next() {
		var category models.BalanceCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan balance category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
