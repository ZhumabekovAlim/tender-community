package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
)

type ChangeRepository struct {
	Db *sql.DB
}

func (r *ChangeRepository) CreateChange(ctx context.Context, change *models.Change) (int, error) {
	query := `
		INSERT INTO changes (transaction_id, amount, description) 
		VALUES (?, ?, ?)
	`
	result, err := r.Db.ExecContext(ctx, query, change.TransactionID, change.Amount, change.Description)
	if err != nil {
		return 0, fmt.Errorf("failed to create change: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	return int(id), nil
}

func (r *ChangeRepository) GetChangeByID(ctx context.Context, id int) (*models.Change, error) {
	query := `
		SELECT id, transaction_id, amount, description, date
		FROM changes
		WHERE id = ?
	`
	var change models.Change
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&change.ID, &change.TransactionID, &change.Amount, &change.Description, &change.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get change: %w", err)
	}
	return &change, nil
}

func (r *ChangeRepository) UpdateChange(ctx context.Context, change *models.Change) (*models.Change, error) {
	updateQuery := `
		UPDATE changes
		SET transaction_id = ?, amount = ?, description = ?
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, updateQuery, change.TransactionID, change.Amount, change.Description, change.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update change: %w", err)
	}

	selectQuery := `
		SELECT id, transaction_id, amount, description, date
		FROM changes
		WHERE id = ?
	`
	var updatedChange models.Change
	err = r.Db.QueryRowContext(ctx, selectQuery, change.ID).Scan(
		&updatedChange.ID, &updatedChange.TransactionID, &updatedChange.Amount, &updatedChange.Description, &updatedChange.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated change: %w", err)
	}

	return &updatedChange, nil
}

func (r *ChangeRepository) DeleteChange(ctx context.Context, id int) error {
	query := `
		DELETE FROM changes
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete change: %w", err)
	}
	return nil
}

func (r *ChangeRepository) GetAllChangesByTransactionID(ctx context.Context, transactionID int) ([]models.Change, error) {
	query := `
		SELECT id, transaction_id, amount, description, date
		FROM changes
		WHERE transaction_id = ?
		ORDER BY date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query changes by transaction_id: %w", err)
	}
	defer rows.Close()

	var changes []models.Change
	for rows.Next() {
		var change models.Change
		if err := rows.Scan(&change.ID, &change.TransactionID, &change.Amount, &change.Description, &change.Date); err != nil {
			return nil, fmt.Errorf("failed to scan change: %w", err)
		}
		changes = append(changes, change)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return changes, nil
}
