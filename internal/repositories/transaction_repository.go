package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tender/internal/models"
)

type TransactionRepository struct {
	Db *sql.DB
}

type MonthlyAmount struct {
	Month  string  `json:"name"`
	Amount float64 `json:"amount"`
}

type YearlyAmounts struct {
	Year   int             `json:"year"`
	Months []MonthlyAmount `json:"months"`
}

type CompanyTotalAmount struct {
	CompanyName string  `json:"name"`
	TotalAmount float64 `json:"total_amount"`
}

type DebtResult struct {
	CompanyID int     `json:"company_id"`
	Debt      float64 `json:"debt"`
}

type DebtResult1 struct {
	TransactionID int     `json:"transaction_id"`
	CompanyID     int     `json:"company_id"`
	Debt          float64 `json:"debt"`
}
type DebtResult2 struct {
	TransactionID int     `json:"transaction_id"`
	Debt          float64 `json:"debt"`
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	// Begin a new database transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Transaction{}, err
	}

	// Insert the transaction
	result, err := tx.ExecContext(ctx, `
    INSERT INTO transactions (transaction_number, type, tender_number, user_id, company_id, organization, amount, total, sell, product_name, completed_date, date, status)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		transaction.TransactionNumber, transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Sell, transaction.ProductName,
		transaction.CompletedDate, transaction.Date, transaction.Status)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return models.Transaction{}, err
	}

	// Get the last inserted transaction ID
	transactionID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return models.Transaction{}, err
	}
	transaction.ID = int(transactionID)

	// Insert the expenses associated with the transaction
	for _, expense := range transaction.Expenses {
		_, err := tx.ExecContext(ctx, `
      INSERT INTO additional_expenses (name, amount, transaction_id)
      VALUES (?, ?, ?)`,
			expense.Name, expense.Amount, transactionID)
		if err != nil {
			tx.Rollback() // Rollback the transaction on error
			return models.Transaction{}, err
		}
	}

	// Retrieve the user and company names
	err = tx.QueryRowContext(ctx, `
    SELECT u.name, c.name 
    FROM transactions t
    JOIN users u ON t.user_id = u.id
    JOIN companies c ON t.company_id = c.id
    WHERE t.id = ?`,
		transaction.ID).Scan(&transaction.UserName, &transaction.CompanyName)
	if err != nil {
		tx.Rollback() // Rollback the transaction on error
		return models.Transaction{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.Transaction{}, err
	}

	return transaction, nil
}

// GetTransactionByID retrieves a transaction by ID from the database along with its expenses.
func (r *TransactionRepository) GetTransactionByID(ctx context.Context, id int) (models.Transaction, error) {
	var transaction models.Transaction

	// Retrieve the transaction
	err := r.Db.QueryRowContext(ctx, `
		SELECT id, transaction_number, type, tender_number, user_id, company_id, organization, amount, total, sell, product_name, completed_date, date, status
		FROM transactions WHERE id = ?`, id).Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
		&transaction.UserID, &transaction.CompanyID, &transaction.Organization, &transaction.Amount,
		&transaction.Total, &transaction.Sell, &transaction.ProductName, &transaction.CompletedDate,
		&transaction.Date, &transaction.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return transaction, models.ErrTransactionNotFound
		}
		return transaction, err
	}

	// Retrieve associated expenses
	rows, err := r.Db.QueryContext(ctx, `
		SELECT id, name, amount, transaction_id
		FROM additional_expenses WHERE transaction_id = ?`, id)
	if err != nil {
		return transaction, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID)
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
		SELECT t.id,t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, c.name, u.name
		FROM transactions t JOIN tender.companies c on c.id = t.company_id JOIN tender.users u on u.id = t.user_id ORDER BY t.date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
			&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
			&transaction.Amount, &transaction.Total, &transaction.Sell, &transaction.ProductName,
			&transaction.CompletedDate, &transaction.Date, &transaction.Status, &transaction.CompanyName, &transaction.UserName)
		if err != nil {
			return nil, err
		}

		// Retrieve associated expenses for each transaction
		expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id
			FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
		if err != nil {
			return nil, err
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID)
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

func (r *TransactionRepository) GetTransactionsByUser(ctx context.Context, userID int) ([]models.Transaction, error) {
	query := `
		SELECT transactions.*, u.name, c.name
		FROM transactions
		JOIN tender.users u ON u.id = transactions.user_id
		JOIN tender.companies c ON c.id = transactions.company_id
		WHERE u.id = ?
		ORDER BY transactions.date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionNumber,
			&transaction.Type,
			&transaction.TenderNumber,
			&transaction.UserID,
			&transaction.CompanyID,
			&transaction.Organization,
			&transaction.Amount,
			&transaction.Total,
			&transaction.Sell,
			&transaction.ProductName,
			&transaction.CompletedDate,
			&transaction.Date,
			&transaction.Status,
			&transaction.UserName,
			&transaction.CompanyName,
		); err != nil {
			return nil, err
		}

		// Retrieve associated expenses for each transaction
		expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id
			FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
		if err != nil {
			return nil, err
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			if err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID); err != nil {
				return nil, err
			}
			expenses = append(expenses, expense)
		}

		transaction.Expenses = expenses
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsByCompany(ctx context.Context, companyID int) ([]models.Transaction, error) {
	query := `
		SELECT transactions.*, u.name AS user_name, c.name AS company_name
		FROM transactions
		JOIN tender.users u ON u.id = transactions.user_id
		JOIN tender.companies c ON c.id = transactions.company_id
		WHERE c.id = ?
		ORDER BY transactions.date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionNumber,
			&transaction.Type,
			&transaction.TenderNumber,
			&transaction.UserID,
			&transaction.CompanyID,
			&transaction.Organization,
			&transaction.Amount,
			&transaction.Total,
			&transaction.Sell,
			&transaction.ProductName,
			&transaction.CompletedDate,
			&transaction.Date,
			&transaction.Status,
			&transaction.UserName,
			&transaction.CompanyName,
		); err != nil {
			return nil, err
		}

		// Retrieve associated expenses for each transaction
		expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id
			FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
		if err != nil {
			return nil, err
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			if err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID); err != nil {
				return nil, err
			}
			expenses = append(expenses, expense)
		}
		transaction.Expenses = expenses

		// Calculate the debt (sell - sum of tranches)
		var totalTranches float64
		err = r.Db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(amount), 0) FROM tranches WHERE transaction_id = ? AND transaction_id IN 
				(SELECT id FROM transactions WHERE company_id = ?)`, transaction.ID, companyID).Scan(&totalTranches)
		if err != nil {
			return nil, err
		}
		transaction.Debt = transaction.Sell - totalTranches

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetTransactionsForUserByCompany(ctx context.Context, userID, companyID int) ([]models.Transaction, error) {
	query := `
		SELECT transactions.*, u.name, c.name
		FROM transactions
		JOIN tender.users u ON u.id = transactions.user_id
		JOIN tender.companies c ON c.id = transactions.company_id
		WHERE c.id = ? AND u.id = ?
		ORDER BY transactions.date DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, companyID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction

		if err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionNumber,
			&transaction.Type,
			&transaction.TenderNumber,
			&transaction.UserID,
			&transaction.CompanyID,
			&transaction.Organization,
			&transaction.Amount,
			&transaction.Total,
			&transaction.Sell,
			&transaction.ProductName,
			&transaction.CompletedDate,
			&transaction.Date,
			&transaction.Status,
			&transaction.UserName,
			&transaction.CompanyName,
		); err != nil {
			return nil, err
		}

		// Retrieve associated expenses for each transaction
		expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id
			FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
		if err != nil {
			return nil, err
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			if err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID); err != nil {
				return nil, err
			}
			expenses = append(expenses, expense)
		}

		transaction.Expenses = expenses
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetAllTransactionsSum(ctx context.Context) (*models.TransactionDebt, error) {
	queryZakup := `
        SELECT COALESCE(SUM(sell), 0) AS total_sum
        FROM tender.transactions
        WHERE status = 3
        AND type = 'Закуп';
    `

	// Execute the query for Zakup
	var totalZakup float64
	err := r.Db.QueryRowContext(ctx, queryZakup).Scan(&totalZakup)
	if err != nil {
		return nil, err
	}

	// Return the result in a struct
	return &models.TransactionDebt{
		Zakup: totalZakup,
	}, nil
}

func (r *TransactionRepository) GetTransactionCountsByUserID(ctx context.Context, userID int) (*models.TransactionCount, error) {

	queryTotal := `
        SELECT COUNT(*) 
        FROM tender.transactions 
        WHERE user_id = ?;
    `

	queryStatus := `
        SELECT COUNT(*) 
        FROM tender.transactions 
        WHERE user_id = ? AND status = ?;
    `

	counts := &models.TransactionCount{}

	err := r.Db.QueryRowContext(ctx, queryTotal, userID).Scan(&counts.TotalTransactions)
	if err != nil {
		return nil, fmt.Errorf("failed to count total transactions: %w", err)
	}

	// Execute queries for each status
	statusCounts := []*int{&counts.Status0, &counts.Status1, &counts.Status2, &counts.Status3}
	for i := 0; i < 4; i++ {
		err = r.Db.QueryRowContext(ctx, queryStatus, userID, i).Scan(statusCounts[i])
		if err != nil {
			return nil, fmt.Errorf("failed to count status %d transactions: %w", i, err)
		}
	}

	return counts, nil
}

func (r *TransactionRepository) GetTransactionsDebtZakup(ctx context.Context, userID int) (*models.TransactionDebt, error) {
	queryZakup := `
        SELECT COALESCE(SUM(total), 0) AS total_sum
        FROM tender.transactions
        WHERE status = 2
        AND user_id = ?
        AND type = 'Закуп';
    `

	// Execute the query for Zakup
	var totalZakup float64
	err := r.Db.QueryRowContext(ctx, queryZakup, userID).Scan(&totalZakup)
	if err != nil {
		return nil, err
	}

	// Return the result in a struct
	return &models.TransactionDebt{
		Zakup: totalZakup,
	}, nil
}

func (r *TransactionRepository) GetTransactionsDebt(ctx context.Context, transactionID int) (*models.TransactionDebtId, error) {
	queryZakup := `
        SELECT total AS total_sum
        FROM tender.transactions
        WHERE status = 2
        AND id = ?
        AND type = 'Закуп';
    `

	// Execute the query for Zakup
	var totalZakup float64
	err := r.Db.QueryRowContext(ctx, queryZakup, transactionID).Scan(&totalZakup)
	if err != nil {
		return nil, err
	}

	// Return the result in a struct
	return &models.TransactionDebtId{
		Debt: totalZakup,
	}, nil
}

// UpdateTransaction updates an existing transaction and its expenses in the database.
func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	// Begin a new database transaction
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		return models.Transaction{}, err
	}

	// Retrieve existing transaction data to preserve non-updated fields
	var existingTransaction models.Transaction
	row := tx.QueryRowContext(ctx, `
		SELECT transaction_number,type, tender_number, user_id, company_id, organization, amount, total, sell, product_name, completed_date, status 
		FROM transactions WHERE id = ?`, transaction.ID)
	err = row.Scan(&existingTransaction.TransactionNumber, &existingTransaction.Type, &existingTransaction.TenderNumber, &existingTransaction.UserID,
		&existingTransaction.CompanyID, &existingTransaction.Organization, &existingTransaction.Amount,
		&existingTransaction.Total, &existingTransaction.Sell, &existingTransaction.ProductName,
		&existingTransaction.CompletedDate, &existingTransaction.Status)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return models.Transaction{}, models.ErrTransactionNotFound
	} else if err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	// Set the values to be updated, preserving existing ones if not provided
	if transaction.TransactionNumber == nil {
		transaction.TransactionNumber = existingTransaction.TransactionNumber
	}
	if transaction.Type == "" {
		transaction.Type = existingTransaction.Type
	}
	if transaction.TenderNumber == nil {
		transaction.TenderNumber = existingTransaction.TenderNumber
	}
	if transaction.UserID == nil {
		transaction.UserID = existingTransaction.UserID
	}
	if transaction.CompanyID == nil {
		transaction.CompanyID = existingTransaction.CompanyID
	}
	if transaction.Organization == nil {
		transaction.Organization = existingTransaction.Organization
	}
	if transaction.Amount == 0 {
		transaction.Amount = existingTransaction.Amount
	}
	if transaction.Total == 0 {
		transaction.Total = existingTransaction.Total
	}
	if transaction.Sell == 0 {
		transaction.Sell = existingTransaction.Sell
	}
	if transaction.ProductName == "" {
		transaction.ProductName = existingTransaction.ProductName
	}
	if transaction.CompletedDate == nil {
		transaction.CompletedDate = existingTransaction.CompletedDate
	}

	// Update the transaction
	result, err := tx.ExecContext(ctx, `
		UPDATE transactions SET transaction_number = ?, type = ?, tender_number = ?, user_id = ?, company_id = ?, 
		organization = ?, amount = ?, total = ?, sell = ?, product_name = ?,  status = ?, completed_date = ? WHERE id = ?`,
		transaction.TransactionNumber, transaction.Type, transaction.TenderNumber, transaction.UserID, transaction.CompanyID,
		transaction.Organization, transaction.Amount, transaction.Total, transaction.Sell,
		transaction.ProductName, transaction.Status, transaction.CompletedDate, transaction.ID)
	if err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	if rowsAffected == 0 {
		tx.Rollback()
		log.Printf("No rows affected for transaction ID: %d", transaction.ID)
		return models.Transaction{}, models.ErrTransactionNotFound
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
	if err != nil {
		tx.Rollback()
		return models.Transaction{}, err
	}

	for _, expense := range transaction.Expenses {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO additional_expenses (name, amount, transaction_id)
			VALUES (?, ?, ?)`,
			expense.Name, expense.Amount, transaction.ID)
		if err != nil {
			tx.Rollback()
			return models.Transaction{}, err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.Transaction{}, err
	}

	// Retrieve the updated transaction data including the user name, company name, and updated expenses
	row = r.Db.QueryRowContext(ctx, `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, u.name, c.name 
		FROM transactions t
		JOIN users u ON t.user_id = u.id
		JOIN companies c ON t.company_id = c.id
		WHERE t.id = ?`, transaction.ID)
	err = row.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber, &transaction.UserID,
		&transaction.CompanyID, &transaction.Organization, &transaction.Amount, &transaction.Total, &transaction.Sell,
		&transaction.ProductName, &transaction.CompletedDate, &transaction.Date, &transaction.Status,
		&transaction.UserName, &transaction.CompanyName)
	if err != nil {
		return models.Transaction{}, err
	}

	// Retrieve the updated expenses
	rows, err := r.Db.QueryContext(ctx, `SELECT id, name, amount, transaction_id FROM additional_expenses WHERE transaction_id = ?`, transaction.ID)
	if err != nil {
		return models.Transaction{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID)
		if err != nil {
			return models.Transaction{}, err
		}
		transaction.Expenses = append(transaction.Expenses, expense)
	}

	return transaction, nil
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

func (r *TransactionRepository) GetMonthlyAmountsByGlobal(ctx context.Context) ([]YearlyAmounts, error) {
	query := `
		SELECT 
			YEAR(date) as year, 
			MONTHNAME(date) as month, 
			SUM(amount) as total_amount 
		FROM 
			transactions 
		GROUP BY 
			year, MONTH(date) 
		ORDER BY 
			year DESC, MONTH(date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var yearlyAmounts []YearlyAmounts
	var currentYear int
	var yearlyData YearlyAmounts
	var monthlyData MonthlyAmount

	for rows.Next() {
		var year int
		var month string
		var totalAmount float64

		if err := rows.Scan(&year, &month, &totalAmount); err != nil {
			return nil, err
		}

		if currentYear != year {
			if currentYear != 0 {
				yearlyAmounts = append(yearlyAmounts, yearlyData)
			}
			yearlyData = YearlyAmounts{
				Year:   year,
				Months: []MonthlyAmount{},
			}
			currentYear = year
		}

		monthlyData = MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		yearlyData.Months = append(yearlyData.Months, monthlyData)
	}

	if currentYear != 0 {
		yearlyAmounts = append(yearlyAmounts, yearlyData)
	}

	return yearlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsByYear(ctx context.Context, year int) ([]MonthlyAmount, error) {
	query := `
		SELECT
			MONTHNAME(date) as month,
			SUM(amount) as total_amount
		FROM
			transactions
		WHERE
			YEAR(date) = ?
		GROUP BY
			MONTH(date)
		ORDER BY
			MONTH(date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyAmounts []MonthlyAmount

	for rows.Next() {
		var month string
		var totalAmount float64

		if err := rows.Scan(&month, &totalAmount); err != nil {
			return nil, err
		}

		monthlyAmount := MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		monthlyAmounts = append(monthlyAmounts, monthlyAmount)
	}

	return monthlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsByCompany(ctx context.Context, companyID int) ([]MonthlyAmount, error) {
	query := `
		SELECT
			MONTHNAME(date) as month,
			SUM(amount) as total_amount
		FROM
			transactions
		WHERE
			company_id = ?
		GROUP BY
			MONTH(date)
		ORDER BY
			MONTH(date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyAmounts []MonthlyAmount

	for rows.Next() {
		var month string
		var totalAmount float64

		if err := rows.Scan(&month, &totalAmount); err != nil {
			return nil, err
		}

		monthlyAmount := MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		monthlyAmounts = append(monthlyAmounts, monthlyAmount)
	}

	return monthlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsByYearAndCompany(ctx context.Context, year int, companyID int) ([]MonthlyAmount, error) {
	query := `
		SELECT
			MONTHNAME(t.date) as month,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		WHERE
			YEAR(t.date) = ? AND c.id = ?
		GROUP BY
			MONTH(t.date)
		ORDER BY
			MONTH(t.date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, year, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyAmounts []MonthlyAmount

	for rows.Next() {
		var month string
		var totalAmount float64

		if err := rows.Scan(&month, &totalAmount); err != nil {
			return nil, err
		}

		monthlyAmount := MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		monthlyAmounts = append(monthlyAmounts, monthlyAmount)
	}

	return monthlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsGroupedByYear(ctx context.Context) ([]YearlyAmounts, error) {
	query := `
		SELECT
			YEAR(date) as year,
			MONTHNAME(date) as month,
			SUM(amount) as total_amount
		FROM
			transactions
		GROUP BY
			year, MONTH(date)
		ORDER BY
			year DESC, MONTH(date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var yearlyAmounts []YearlyAmounts
	var currentYear int
	var yearlyData YearlyAmounts
	var monthlyData MonthlyAmount

	for rows.Next() {
		var year int
		var month string
		var totalAmount float64

		if err := rows.Scan(&year, &month, &totalAmount); err != nil {
			return nil, err
		}

		if currentYear != year {
			if currentYear != 0 {
				yearlyAmounts = append(yearlyAmounts, yearlyData)
			}
			yearlyData = YearlyAmounts{
				Year:   year,
				Months: []MonthlyAmount{},
			}
			currentYear = year
		}

		monthlyData = MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		yearlyData.Months = append(yearlyData.Months, monthlyData)
	}

	if currentYear != 0 {
		yearlyAmounts = append(yearlyAmounts, yearlyData)
	}

	return yearlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsGroupedByYearForUser(ctx context.Context, userID int) ([]YearlyAmounts, error) {
	query := `
		SELECT
			YEAR(t.date) as year,
			MONTHNAME(t.date) as month,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			users u on u.id = t.user_id
		WHERE 
			u.id = ?
		GROUP BY
			year, MONTH(t.date)
		ORDER BY
			year DESC, MONTH(t.date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var yearlyAmounts []YearlyAmounts
	var currentYear int
	var yearlyData YearlyAmounts
	var monthlyData MonthlyAmount

	for rows.Next() {
		var year int
		var month string
		var totalAmount float64

		if err := rows.Scan(&year, &month, &totalAmount); err != nil {
			return nil, err
		}

		if currentYear != year {
			if currentYear != 0 {
				yearlyAmounts = append(yearlyAmounts, yearlyData)
			}
			yearlyData = YearlyAmounts{
				Year:   year,
				Months: []MonthlyAmount{},
			}
			currentYear = year
		}

		monthlyData = MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		yearlyData.Months = append(yearlyData.Months, monthlyData)
	}

	if currentYear != 0 {
		yearlyAmounts = append(yearlyAmounts, yearlyData)
	}

	return yearlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsForUserByYear(ctx context.Context, userID int, year int) ([]MonthlyAmount, error) {
	query := `
		SELECT
			MONTHNAME(t.date) as month,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			users u on u.id = t.user_id
		WHERE 
			u.id = ? AND YEAR(t.date) = ?
		GROUP BY
			MONTH(t.date)
		ORDER BY
			MONTH(t.date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyAmounts []MonthlyAmount

	for rows.Next() {
		var month string
		var totalAmount float64

		if err := rows.Scan(&month, &totalAmount); err != nil {
			return nil, err
		}

		monthlyAmount := MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		monthlyAmounts = append(monthlyAmounts, monthlyAmount)
	}

	return monthlyAmounts, nil
}

func (r *TransactionRepository) GetMonthlyAmountsForUserByYearAndCompany(ctx context.Context, userID int, year int, companyID int) ([]MonthlyAmount, error) {
	query := `
		SELECT
			MONTHNAME(t.date) as month,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			users u on u.id = t.user_id
		JOIN
			companies c on c.id = t.company_id
		WHERE 
			u.id = ? AND YEAR(t.date) = ? AND c.id = ?
		GROUP BY
			MONTH(t.date)
		ORDER BY
			MONTH(t.date) DESC
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, year, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyAmounts []MonthlyAmount

	for rows.Next() {
		var month string
		var totalAmount float64

		if err := rows.Scan(&month, &totalAmount); err != nil {
			return nil, err
		}

		monthlyAmount := MonthlyAmount{
			Month:  month,
			Amount: totalAmount,
		}
		monthlyAmounts = append(monthlyAmounts, monthlyAmount)
	}

	return monthlyAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountGroupedByCompany(ctx context.Context) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForYear(ctx context.Context, year int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		WHERE
			YEAR(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForMonth(ctx context.Context, month int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN
			companies c on c.id = t.company_id
		WHERE
			MONTH(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForYearAndMonth(ctx context.Context, year int, month int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		WHERE
			YEAR(t.date) = ? AND MONTH(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountGroupedByCompanyForUsers(ctx context.Context) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForUser(ctx context.Context, userID int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		JOIN
			users u on u.id = t.user_id
		WHERE 
			u.id = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForUserAndMonth(ctx context.Context, userID int, month int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN
			companies c on c.id = t.company_id
		JOIN
			users u on u.id = t.user_id
		WHERE
			u.id = ? AND MONTH(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForUserAndYear(ctx context.Context, userID int, year int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN
			companies c on c.id = t.company_id
		JOIN
			users u on u.id = t.user_id
		WHERE 
			u.id = ? AND YEAR(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) GetTotalAmountByCompanyForUserYearAndMonth(ctx context.Context, userID int, year int, month int) ([]CompanyTotalAmount, error) {
	query := `
		SELECT
			c.name,
			SUM(t.amount) as total_amount
		FROM
			transactions t
		JOIN 
			companies c on c.id = t.company_id
		JOIN
			users u on u.id = t.user_id
		WHERE 
			u.id = ? AND YEAR(t.date) = ? AND MONTH(t.date) = ?
		GROUP BY
			c.id
		ORDER BY
			c.id
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, year, month)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalAmounts []CompanyTotalAmount

	for rows.Next() {
		var companyName string
		var totalAmount float64

		if err := rows.Scan(&companyName, &totalAmount); err != nil {
			return nil, err
		}

		companyTotal := CompanyTotalAmount{
			CompanyName: companyName,
			TotalAmount: totalAmount,
		}
		totalAmounts = append(totalAmounts, companyTotal)
	}

	return totalAmounts, nil
}

func (r *TransactionRepository) FindAllTransactionsByUserIDAndStatus(ctx context.Context, userID, status int) ([]models.Transaction, error) {
	query := `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, 
		       t.organization, t.amount, t.total, t.sell, t.product_name, t.completed_date, 
		       t.date, t.status, u.name as username, c.name as companyname
		FROM transactions t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN companies c ON t.company_id = c.id
		WHERE t.user_id = ? AND t.status = ?;
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TransactionNumber, &t.Type, &t.TenderNumber, &t.UserID, &t.CompanyID, &t.Organization,
			&t.Amount, &t.Total, &t.Sell, &t.ProductName, &t.CompletedDate, &t.Date, &t.Status, &t.UserName, &t.CompanyName); err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}

		// Retrieve associated expenses for each transaction
		expenseQuery := `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses 
			WHERE transaction_id = ?
		`
		expenseRows, err := r.Db.QueryContext(ctx, expenseQuery, t.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to query expenses for transaction %d: %w", t.ID, err)
		}
		defer expenseRows.Close()

		var expenses []models.Expense
		for expenseRows.Next() {
			var expense models.Expense
			if err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date); err != nil {
				return nil, fmt.Errorf("failed to scan expense: %w", err)
			}
			expenses = append(expenses, expense)
		}
		expenseRows.Close() // Close this set after processing expenses for each transaction

		// Assign the expenses to the transaction
		t.Expenses = expenses
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *TransactionRepository) FindAllTendersByUserIDAndStatus(ctx context.Context, userID, status int) ([]models.Tender, error) {
	query := `
		SELECT t.id, t.type, t.tender_number, t.user_id, t.company_id, t.organization, t.total, 
		       t.commission, t.completed_date, t.date, t.status, u.name as username, c.name as companyname
		FROM tenders t
		LEFT JOIN users u ON t.user_id = u.id
		LEFT JOIN companies c ON t.company_id = c.id
		WHERE t.user_id = ? AND t.status = ?;
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query tenders: %w", err)
	}
	defer rows.Close()

	var tenders []models.Tender
	for rows.Next() {
		var t models.Tender
		if err := rows.Scan(&t.ID, &t.Type, &t.TenderNumber, &t.UserID, &t.CompanyID, &t.Organization, &t.Total, &t.Commission, &t.CompletedDate, &t.Date, &t.Status, &t.UserName, &t.CompanyName); err != nil {
			return nil, fmt.Errorf("failed to scan tender: %w", err)
		}
		tenders = append(tenders, t)
	}
	return tenders, nil
}

func (r *TransactionRepository) FindAllExtraTransactionsByUserIDAndStatus(ctx context.Context, userID, status int) ([]models.ExtraTransaction, error) {
	query := `
		SELECT e.id, e.user_id, e.description, e.total, e.date, e.status, u.name as username
		FROM extra_transactions e
		LEFT JOIN users u ON e.user_id = u.id
		WHERE e.user_id = ? AND e.status = ?;
	`

	rows, err := r.Db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to query extra transactions: %w", err)
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var et models.ExtraTransaction
		if err := rows.Scan(&et.ID, &et.UserID, &et.Description, &et.Total, &et.Date, &et.Status, &et.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan extra transaction: %w", err)
		}
		extraTransactions = append(extraTransactions, et)
	}
	return extraTransactions, nil
}

func (r *TransactionRepository) GetAllTransactionsByDateRange(ctx context.Context, startDate, endDate string, userId int) ([]models.Transaction, error) {
	if userId > 0 {
		if userId == 1 {
			query := `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
		       t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, 
		       c.name AS companyname, u.name AS username
		FROM transactions t 
		JOIN companies c ON c.id = t.company_id 
		JOIN users u ON u.id = t.user_id 
		WHERE t.completed_date BETWEEN ? AND ? 
		ORDER BY t.completed_date DESC
	`
			rows, err := r.Db.QueryContext(ctx, query, startDate, endDate)
			if err != nil {
				return nil, fmt.Errorf("failed to query transactions: %w", err)
			}

			defer rows.Close()

			var transactions []models.Transaction
			for rows.Next() {
				var transaction models.Transaction
				err := rows.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
					&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
					&transaction.Amount, &transaction.Total, &transaction.Sell, &transaction.ProductName,
					&transaction.CompletedDate, &transaction.Date, &transaction.Status, &transaction.CompanyName, &transaction.UserName)
				if err != nil {
					return nil, fmt.Errorf("failed to scan transaction: %w", err)
				}

				// Retrieve associated expenses for each transaction
				expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses 
			WHERE transaction_id = ?`, transaction.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to query expenses for transaction %d: %w", transaction.ID, err)
				}
				defer expenseRows.Close()

				var expenses []models.Expense
				for expenseRows.Next() {
					var expense models.Expense
					err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
					if err != nil {
						return nil, fmt.Errorf("failed to scan expense: %w", err)
					}
					expenses = append(expenses, expense)
				}
				expenseRows.Close() // Close after processing expenses for each transaction

				transaction.Expenses = expenses
				transactions = append(transactions, transaction)
			}

			// Check for errors during row iteration
			if err := rows.Err(); err != nil {
				return nil, err
			}

			return transactions, nil
		} else {
			query := `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
		       t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, 
		       c.name AS companyname, u.name AS username
		FROM transactions t 
		JOIN companies c ON c.id = t.company_id 
		JOIN users u ON u.id = t.user_id 
		WHERE t.user_id = ? AND t.completed_date BETWEEN ? AND ? 
		ORDER BY t.completed_date DESC
	`
			rows, err := r.Db.QueryContext(ctx, query, userId, startDate, endDate)
			if err != nil {
				return nil, fmt.Errorf("failed to query transactions: %w", err)
			}

			defer rows.Close()

			var transactions []models.Transaction
			for rows.Next() {
				var transaction models.Transaction
				err := rows.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
					&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
					&transaction.Amount, &transaction.Total, &transaction.Sell, &transaction.ProductName,
					&transaction.CompletedDate, &transaction.Date, &transaction.Status, &transaction.CompanyName, &transaction.UserName)
				if err != nil {
					return nil, fmt.Errorf("failed to scan transaction: %w", err)
				}

				// Retrieve associated expenses for each transaction
				expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses 
			WHERE transaction_id = ?`, transaction.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to query expenses for transaction %d: %w", transaction.ID, err)
				}
				defer expenseRows.Close()

				var expenses []models.Expense
				for expenseRows.Next() {
					var expense models.Expense
					err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
					if err != nil {
						return nil, fmt.Errorf("failed to scan expense: %w", err)
					}
					expenses = append(expenses, expense)
				}
				expenseRows.Close() // Close after processing expenses for each transaction

				transaction.Expenses = expenses
				transactions = append(transactions, transaction)
			}

			// Check for errors during row iteration
			if err := rows.Err(); err != nil {
				return nil, err
			}

			return transactions, nil
		}
	}
	return []models.Transaction{}, nil
}

func (r *TransactionRepository) GetAllTransactionsByDateRangeCompany(ctx context.Context, startDate, endDate string, userId, companyId int) ([]models.Transaction, error) {
	if userId > 0 {
		if userId == 1 {
			query := `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
		       t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, 
		       c.name AS companyname, u.name AS username
		FROM transactions t 
		JOIN companies c ON c.id = t.company_id 
		JOIN users u ON u.id = t.user_id 
		WHERE t.company_id = ? AND t.completed_date BETWEEN ? AND ? 
		ORDER BY t.completed_date DESC
	`
			rows, err := r.Db.QueryContext(ctx, query, companyId, startDate, endDate)
			if err != nil {
				return nil, fmt.Errorf("failed to query transactions: %w", err)
			}

			defer rows.Close()

			var transactions []models.Transaction
			for rows.Next() {
				var transaction models.Transaction
				err := rows.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
					&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
					&transaction.Amount, &transaction.Total, &transaction.Sell, &transaction.ProductName,
					&transaction.CompletedDate, &transaction.Date, &transaction.Status, &transaction.CompanyName, &transaction.UserName)
				if err != nil {
					return nil, fmt.Errorf("failed to scan transaction: %w", err)
				}

				// Retrieve associated expenses for each transaction
				expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses 
			WHERE transaction_id = ?`, transaction.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to query expenses for transaction %d: %w", transaction.ID, err)
				}
				defer expenseRows.Close()

				var expenses []models.Expense
				for expenseRows.Next() {
					var expense models.Expense
					err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
					if err != nil {
						return nil, fmt.Errorf("failed to scan expense: %w", err)
					}
					expenses = append(expenses, expense)
				}
				expenseRows.Close() // Close after processing expenses for each transaction

				transaction.Expenses = expenses
				transactions = append(transactions, transaction)
			}

			// Check for errors during row iteration
			if err := rows.Err(); err != nil {
				return nil, err
			}

			return transactions, nil
		} else {
			query := `
		SELECT t.id, t.transaction_number, t.type, t.tender_number, t.user_id, t.company_id, t.organization, 
		       t.amount, t.total, t.sell, t.product_name, t.completed_date, t.date, t.status, 
		       c.name AS companyname, u.name AS username
		FROM transactions t 
		JOIN companies c ON c.id = t.company_id 
		JOIN users u ON u.id = t.user_id 
		WHERE t.user_id = ? AND t.company_id = ? AND t.completed_date BETWEEN ? AND ? 
		ORDER BY t.completed_date DESC
	`
			rows, err := r.Db.QueryContext(ctx, query, userId, companyId, startDate, endDate)
			if err != nil {
				return nil, fmt.Errorf("failed to query transactions: %w", err)
			}

			defer rows.Close()

			var transactions []models.Transaction
			for rows.Next() {
				var transaction models.Transaction
				err := rows.Scan(&transaction.ID, &transaction.TransactionNumber, &transaction.Type, &transaction.TenderNumber,
					&transaction.UserID, &transaction.CompanyID, &transaction.Organization,
					&transaction.Amount, &transaction.Total, &transaction.Sell, &transaction.ProductName,
					&transaction.CompletedDate, &transaction.Date, &transaction.Status, &transaction.CompanyName, &transaction.UserName)
				if err != nil {
					return nil, fmt.Errorf("failed to scan transaction: %w", err)
				}

				// Retrieve associated expenses for each transaction
				expenseRows, err := r.Db.QueryContext(ctx, `
			SELECT id, name, amount, transaction_id, date
			FROM additional_expenses 
			WHERE transaction_id = ?`, transaction.ID)
				if err != nil {
					return nil, fmt.Errorf("failed to query expenses for transaction %d: %w", transaction.ID, err)
				}
				defer expenseRows.Close()

				var expenses []models.Expense
				for expenseRows.Next() {
					var expense models.Expense
					err := expenseRows.Scan(&expense.ID, &expense.Name, &expense.Amount, &expense.TransactionID, &expense.Date)
					if err != nil {
						return nil, fmt.Errorf("failed to scan expense: %w", err)
					}
					expenses = append(expenses, expense)
				}
				expenseRows.Close() // Close after processing expenses for each transaction

				transaction.Expenses = expenses
				transactions = append(transactions, transaction)
			}

			// Check for errors during row iteration
			if err := rows.Err(); err != nil {
				return nil, err
			}

			return transactions, nil
		}
	}
	return []models.Transaction{}, nil
}

// GetCompanyDebt retrieves the total debt for each company based on transactions with status=2.
func (r *TransactionRepository) GetCompanyDebt(ctx context.Context) ([]DebtResult, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT 
			t.company_id, 
			SUM(t.sell - IFNULL(tr.total_tranche_amount, 0)) AS debt
		FROM 
			transactions t
		LEFT JOIN (
			SELECT 
				transaction_id, 
				SUM(amount) AS total_tranche_amount
			FROM 
				tranches
			GROUP BY 
				transaction_id
		) tr ON t.id = tr.transaction_id
		WHERE 
			t.status = 2
		GROUP BY 
			t.company_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DebtResult
	for rows.Next() {
		var result DebtResult
		err := rows.Scan(&result.CompanyID, &result.Debt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetCompanyDebt retrieves the total debt for each company based on transactions with status=2.
func (r *TransactionRepository) GetCompanyDebtId(ctx context.Context, id int) ([]DebtResult2, error) {
	rows, err := r.Db.QueryContext(ctx, `
SELECT
    t.id,
    SUM(t.sell - IFNULL(tr.total_tranche_amount, 0)) AS debt
FROM
    transactions t
        LEFT JOIN (
        SELECT
            transaction_id,
            SUM(amount) AS total_tranche_amount
        FROM
            tranches
        GROUP BY
            transaction_id
    ) tr ON t.id = tr.transaction_id
WHERE transaction_id = ?
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DebtResult2
	for rows.Next() {
		var result DebtResult2
		err := rows.Scan(&result.TransactionID, &result.Debt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *TransactionRepository) GetCompanyDebtById(ctx context.Context) ([]DebtResult1, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT
			t.id AS transaction_id,
			t.company_id,
			t.sell - IFNULL(SUM(tr.amount), 0) AS debt
		FROM
			transactions t
				LEFT JOIN
			tranches tr ON t.id = tr.transaction_id
		WHERE
			t.status = 2 OR t.status = 3
		GROUP BY
			t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []DebtResult1
	for rows.Next() {
		var result DebtResult1
		err := rows.Scan(&result.TransactionID, &result.CompanyID, &result.Debt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	// Check for errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
