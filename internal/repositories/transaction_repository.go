package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type TransactionRepository struct {
	Db *sql.DB
}

// CreateTransaction inserts a new transaction into the database.
func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	_, err := r.Db.ExecContext(ctx, `
		INSERT INTO transactions (type, tender_number, user_id, company_id, organization, amount, total, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Status)
	return err
}

// GetTransactionByID retrieves a transaction by ID from the database.
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.Db.QueryRowContext(ctx, `
		SELECT id, type, tender_number, user_id, company_id, organization, amount, total, date, status
		FROM transactions WHERE id = ?`, id).Scan(&transaction.ID, &transaction.Type, &transaction.TenderNumber,
		&transaction.UserID, &transaction.CompanyID, &transaction.Organization, &transaction.Amount,
		&transaction.Total, &transaction.Date, &transaction.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return transaction, models.ErrTransactionNotFound
		}
		return transaction, err
	}

	return transaction, nil
}

// GetAllTransactions retrieves all transactions from the database.
func (r *TransactionRepository) GetAllTransactions(ctx context.Context) ([]models.Transaction, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT id, type, tender_number, user_id, company_id, organization, amount, total, date, status
		FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.Type, &transaction.TenderNumber,
			&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
			&transaction.Amount, &transaction.Total, &transaction.Date, &transaction.Status)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction updates an existing transaction in the database.
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {
	result, err := r.Db.ExecContext(ctx, `
		UPDATE transactions SET type = ?, tender_number = ?, user_id = ?, company_id = ?, 
		organization = ?, amount = ?, total = ?, status = ? WHERE id = ?`,
		transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Status, transaction.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrTransactionNotFound
	}

	return nil
}

// DeleteTransaction removes a transaction from the database by ID.
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrTransactionNotFound
	}

	return nil
}
