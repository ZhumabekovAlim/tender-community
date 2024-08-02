package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type TransactionRepository struct {
	Db *sql.DB
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) error {
	// Begin a new database transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Insert the transaction
	result, err := tx.ExecContext(ctx, `
		INSERT INTO transactions (type, tender_number, user_id, company_id, organization, amount, total, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Status)
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		} // Rollback the transaction on error
		return err
	}

	// Get the last inserted transaction ID
	transactionID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return err
	}

	// Insert the expenses associated with the transaction
	for _, expense := range transaction.Expenses {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO additional_expenses (name, amount, transaction_id)
			VALUES (?, ?, ?)`,
			expense.Name, expense.Amount, transactionID)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetTransactionByID retrieves a transaction by ID from the database along with its expenses.
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction

	// Retrieve the transaction
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

	// Retrieve associated expenses
	rows, err := r.Db.QueryContext(ctx, `
		SELECT id, name, amount, transaction_id, date
		FROM additional_expenses WHERE transaction_id = ?`, id)
	if err != nil {
		return transaction, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
		if err != nil {
			return transaction, err
		}
		expenses = append(expenses, expense)
	}

	transaction.Expenses = expenses

	return transaction, nil
}

// GetAllTransactions retrieves all transactions from the database along with their expenses.
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

		// Retrieve associated expenses for each transaction
		expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
		if err != nil {
			return nil, err
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
			if err != nil {
				return nil, err
			}
			expenses = append(expenses, expense)
		}

		transaction.Expenses = expenses
		transactions = append(transactions, transaction)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

// UpdateTransaction updates an existing transaction and its expenses in the database.
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction models.Transaction) error {
	// Begin a new database transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Update the transaction
	result, err := tx.ExecContext(ctx, `
		UPDATE transactions SET type = ?, tender_number = ?, user_id = ?, company_id = ?, 
		organization = ?, amount = ?, total = ?, status = ? WHERE id = ?`,
		transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Status, transaction.ID)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return models.ErrTransactionNotFound
	}

	// Delete existing expenses
	_, err = tx.ExecContext(ctx, `DELETE FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert the updated expenses
	for _, expense := range transaction.Expenses {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO additional_expenses (name, amount, transaction_id)
			VALUES (?, ?, ?)`,
			expense.Name, expense.Amount, transaction.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// DeleteTransaction removes a transaction and its expenses from the database by ID.
func (r *TransactionRepository) DeleteTransaction(ctx context.Context, id int) error {
	// Begin a new database transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Delete expenses first due to foreign key constraints
	_, err = tx.ExecContext(ctx, `DELETE FROM additional_expenses WHERE transaction_id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete the transaction
	result, err := tx.ExecContext(ctx, `DELETE FROM transactions WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return models.ErrTransactionNotFound
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
