package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"tender/internal/models"
)

type PersonalDebtRepository struct {
	Db *sql.DB
}

func (r *PersonalDebtRepository) CreatePersonalDebt(ctx context.Context, debt *models.PersonalDebt) (int, error) {
	query := `
		INSERT INTO personal_debts (name, amount, type, get_date, return_date, status) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.Db.ExecContext(ctx, query, debt.Name, debt.Amount, debt.Type, debt.GetDate, debt.ReturnDate, debt.Status)
	if err != nil {
		return 0, fmt.Errorf("failed to create personal debt: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last insert ID: %w", err)
	}
	return int(id), nil
}

func (r *PersonalDebtRepository) GetPersonalDebtByID(ctx context.Context, id int) (*models.PersonalDebt, error) {
	query := `
		SELECT id, name, amount, type, get_date, return_date, status, created_at, updated_at
		FROM personal_debts
		WHERE id = ?
	`
	var debt models.PersonalDebt
	err := r.Db.QueryRowContext(ctx, query, id).Scan(
		&debt.ID, &debt.Name, &debt.Amount, &debt.Type, &debt.GetDate, &debt.ReturnDate, &debt.Status, &debt.CreatedAt, &debt.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get personal debt: %w", err)
	}
	return &debt, nil
}

func (r *PersonalDebtRepository) UpdatePersonalDebt(ctx context.Context, debt *models.PersonalDebt) (*models.PersonalDebt, error) {
	query := `
		UPDATE personal_debts
		SET name = ?, amount = ?, type = ?, get_date = ?, return_date = ?, status = ?
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, debt.Name, debt.Amount, debt.Type, debt.GetDate, debt.ReturnDate, debt.Status, debt.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update personal debt: %w", err)
	}

	// Retrieve the updated personal debt
	return r.GetPersonalDebtByID(ctx, debt.ID)
}

func (r *PersonalDebtRepository) DeletePersonalDebt(ctx context.Context, id int) error {
	query := `
		DELETE FROM personal_debts
		WHERE id = ?
	`
	_, err := r.Db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete personal debt: %w", err)
	}
	return nil
}

func (r *PersonalDebtRepository) GetAllPersonalDebts(ctx context.Context) ([]models.PersonalDebt, error) {
	query := `
		SELECT id, name, amount, type, get_date, return_date, status, created_at, updated_at
		FROM personal_debts
		ORDER BY created_at DESC
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal debts: %w", err)
	}
	defer rows.Close()

	var debts []models.PersonalDebt
	for rows.Next() {
		var debt models.PersonalDebt
		if err := rows.Scan(&debt.ID, &debt.Name, &debt.Amount, &debt.Type, &debt.GetDate, &debt.ReturnDate, &debt.Status, &debt.CreatedAt, &debt.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan personal debt: %w", err)
		}
		debts = append(debts, debt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return debts, nil
}

func (r *PersonalDebtRepository) GetAllPersonalDebtsByStatus(ctx context.Context, status int) ([]models.PersonalDebt, error) {
	query := `
		SELECT id, name, amount, type, get_date, return_date, status, created_at, updated_at
		FROM personal_debts WHERE status = ?
		ORDER BY created_at DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query personal debts: %w", err)
	}
	defer rows.Close()

	var debts []models.PersonalDebt
	for rows.Next() {
		var debt models.PersonalDebt
		if err := rows.Scan(&debt.ID, &debt.Name, &debt.Amount, &debt.Type, &debt.GetDate, &debt.ReturnDate, &debt.Status, &debt.CreatedAt, &debt.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan personal debt: %w", err)
		}
		debts = append(debts, debt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return debts, nil
}
