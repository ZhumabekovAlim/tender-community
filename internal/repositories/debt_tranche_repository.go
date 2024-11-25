package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
)

type DebtTrancheRepository struct {
	Db *sql.DB
}

// Create a new tranche
func (r *DebtTrancheRepository) CreateDebtTranche(ctx context.Context, tranche *models.DebtTranche) (int, error) {
	query := `
		INSERT INTO debt_tranches (debt_id, amount, description, date) 
		VALUES (?, ?, ?, ?)
	`
	result, err := r.Db.ExecContext(ctx, query, tranche.DebtID, tranche.Amount, tranche.Description, tranche.Date)
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
func (r *DebtTrancheRepository) GetDebtTrancheByID(ctx context.Context, id int) (*models.DebtTranche, error) {
	query := `
		SELECT id, debt_id, amount, description, date
		FROM debt_tranches
		WHERE id = ?
	`
	var tranche models.DebtTranche
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&tranche.ID, &tranche.DebtID, &tranche.Amount, &tranche.Description, &tranche.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("failed to get tranche: %w", err)
	}
	return &tranche, nil
}

// Update an existing tranche
func (r *DebtTrancheRepository) UpdateDebtTranche(ctx context.Context, tranche *models.DebtTranche) (*models.DebtTranche, error) {
	// Update the tranche
	updateQuery := `
		UPDATE debt_tranches
		SET debt_id = ?, amount = ?, description = ?, date = ?
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, updateQuery, tranche.DebtID, tranche.Amount, tranche.Description, tranche.Date, tranche.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update tranche: %w", err)
	}

	// Retrieve the updated tranche
	selectQuery := `
		SELECT id, debt_id, amount, description, date
		FROM debt_tranches
		WHERE id = ?
	`

	var updatedTranche models.DebtTranche
	err = r.Db.QueryRowContext(ctx, selectQuery, tranche.ID).Scan(
		&updatedTranche.ID, &updatedTranche.DebtID, &updatedTranche.Amount, &updatedTranche.Description, &updatedTranche.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated tranche: %w", err)
	}

	return &updatedTranche, nil
}

// Delete a tranche by ID
func (r *DebtTrancheRepository) DeleteDebtTranche(ctx context.Context, id int) error {
	query := `
		DELETE FROM debt_tranches
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tranche: %w", err)
	}
	return nil
}

func (r *DebtTrancheRepository) GetAllDebtTranchesByDebtID(ctx context.Context, DebtID int) ([]models.DebtTranche, error) {
	query := `
		SELECT id, debt_id, amount, description, date
		FROM debt_tranches
		WHERE debt_id = ?
		ORDER BY date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, DebtID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tranches by debt_id: %w", err)
	}
	defer rows.Close()

	var tranches []models.DebtTranche
	for rows.Next() {
		var tranche models.DebtTranche
		if err := rows.Scan(&tranche.ID, &tranche.DebtID, &tranche.Amount, &tranche.Description, &tranche.Date); err != nil {
			return nil, fmt.Errorf("failed to scan tranche: %w", err)
		}
		tranches = append(tranches, tranche)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tranches, nil
}
