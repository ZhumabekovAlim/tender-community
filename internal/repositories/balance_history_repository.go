package repositories

import (
	"context"
	"database/sql"
	"errors"
	"tender/internal/models"
)

type BalanceHistoryRepository struct {
	Db *sql.DB
}

// CreateBalanceHistory inserts a new balance history record into the database.
func (r *BalanceHistoryRepository) CreateBalanceHistory(ctx context.Context, history models.BalanceHistory) (models.BalanceHistory, error) {
	result, err := r.Db.ExecContext(ctx, "INSERT INTO balance_history (amount, description, user_id) VALUES (?, ?, ?)",
		history.Amount, history.Description, history.UserID)
	if err != nil {
		return models.BalanceHistory{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.BalanceHistory{}, err
	}

	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, description, user_id, created_at, updated_at FROM balance_history WHERE id = ?", id)
	if err != nil {
		return models.BalanceHistory{}, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&history.ID, &history.Amount, &history.Description, &history.UserID, &history.CreatedAt, &history.UpdatedAt)
		if err != nil {
			return models.BalanceHistory{}, err
		}
	} else {
		return models.BalanceHistory{}, errors.New("no rows found")
	}

	return history, nil
}

func (r *BalanceHistoryRepository) DeleteBalanceHistory(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM balance_history WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrBalanceHistoryNotFound
	}

	return nil
}

// UpdateBalanceHistory updates an existing balance history record in the database.
func (r *BalanceHistoryRepository) UpdateBalanceHistory(ctx context.Context, history models.BalanceHistory) (models.BalanceHistory, error) {
	_, err := r.Db.ExecContext(ctx, "UPDATE balance_history SET amount = ?, description = ?, user_id = ? WHERE id = ?",
		history.Amount, history.Description, history.UserID, history.ID)
	if err != nil {
		return models.BalanceHistory{}, err
	}

	return history, nil
}

// GetAllBalanceHistories retrieves all balance history records from the database.
func (r *BalanceHistoryRepository) GetBalanceHistoryByUserID(ctx context.Context, id int) ([]models.BalanceHistory, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, description, user_id, created_at, updated_at FROM balance_history WHERE user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var histories []models.BalanceHistory
	for rows.Next() {
		var history models.BalanceHistory
		err := rows.Scan(&history.ID, &history.Amount, &history.Description, &history.UserID, &history.CreatedAt, &history.UpdatedAt)
		if err != nil {
			return nil, err
		}
		histories = append(histories, history)
	}

	return histories, nil
}
