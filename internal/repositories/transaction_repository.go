package repositories

import (
	"context"
	"database/sql"
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
		SELECT transactions.*, u.name, c.name
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

	queryTender := `
        SELECT COALESCE(SUM(total), 0) AS total_sum
        FROM tender.transactions
        WHERE status = 2
        AND user_id = ?
        AND (type = 'ГОПП' OR type = 'ГОИК');
    `

	// Execute the query for Tender
	var totalTender float64
	err = r.Db.QueryRowContext(ctx, queryTender, userID).Scan(&totalTender)
	if err != nil {
		return nil, err
	}

	// Return the result in a struct
	return &models.TransactionDebt{
		Zakup:  totalZakup,
		Tender: totalTender,
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
	if transaction.UserID == 0 {
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
