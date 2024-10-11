package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
)

type TrancheRepository struct {
	Db *sql.DB
}

// Create a new tranche
func (r *TrancheRepository) CreateTranche(ctx context.Context, tranche *models.Tranche) (int, error) {
	query := `
		INSERT INTO tranches (transaction_id, amount, description) 
		VALUES (?, ?, ?)
	`
	result, err := r.Db.ExecContext(ctx, query, tranche.TransactionID, tranche.Amount, tranche.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to create tranche: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	return int(id), nil
}

// Get a tranche by ID
func (r *TrancheRepository) GetTrancheByID(ctx context.Context, id int) (*models.Tranche, error) {
	query := `
		SELECT id, transaction_id, amount, description, date
		FROM tranches
		WHERE id = ?
	`
	var tranche models.Tranche
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&tranche.ID, &tranche.TransactionID, &tranche.Amount, &tranche.Description, &tranche.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to get tranche: %w", err)
	}
	return &tranche, nil
}

// Update an existing tranche
func (r *TrancheRepository) UpdateTranche(ctx context.Context, tranche *models.Tranche) (*models.Tranche, error) {
	// Update the tranche
	updateQuery := `
		UPDATE tranches
		SET transaction_id = ?, amount = ?, description = ?
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, updateQuery, tranche.TransactionID, tranche.Amount, tranche.Description, tranche.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update tranche: %w", err)
	}

	// Retrieve the updated tranche
	selectQuery := `
		SELECT id, transaction_id, amount, description, date
		FROM tranches
		WHERE id = ?
	`
	var updatedTranche models.Tranche
	err = r.Db.QueryRowContext(ctx, selectQuery, tranche.ID).Scan(
		&updatedTranche.ID, &updatedTranche.TransactionID, &updatedTranche.Amount, &updatedTranche.Description, &updatedTranche.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated tranche: %w", err)
	}

	return &updatedTranche, nil
}

// Delete a tranche by ID
func (r *TrancheRepository) DeleteTranche(ctx context.Context, id int) error {
	query := `
		DELETE FROM tranches
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tranche: %w", err)
	}
	return nil
}

func (r *TrancheRepository) GetAllTranchesByTransactionID(ctx context.Context, transactionID int) ([]models.Tranche, error) {
	query := `
		SELECT id, transaction_id, amount, description, date
		FROM tranches
		WHERE transaction_id = ?
		ORDER BY date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tranches by transaction_id: %w", err)
	}
	defer rows.Close()

	var tranches []models.Tranche
	for rows.Next() {
		var tranche models.Tranche
		if err := rows.Scan(&tranche.ID, &tranche.TransactionID, &tranche.Amount, &tranche.Description, &tranche.Date); err != nil {
			return nil, fmt.Errorf("failed to scan tranche: %w", err)
		}
		tranches = append(tranches, tranche)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tranches, nil
}
