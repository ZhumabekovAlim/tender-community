package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
	"time"
)

// PersonalExpenseRepository handles database operations for personal expenses.
type PersonalExpenseRepository struct {
	Db *sql.DB
}

// TODO: add date
// CreatePersonalExpense inserts a new expense into the database.
func (r *PersonalExpenseRepository) CreatePersonalExpense(ctx context.Context, expense models.PersonalExpense) (int, error) {
	result, err := r.Db.ExecContext(ctx, "INSERT INTO personal_expenses (amount, reason, description, category_id, date) VALUES (?, ?, ?, ?, ?)", expense.Amount, expense.Reason, expense.Description, expense.CategoryID, expense.Date)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetPersonalExpenseByID retrieves an expense by ID from the database.
func (r *PersonalExpenseRepository) GetPersonalExpenseByID(ctx context.Context, id int) (models.PersonalExpense, error) {
	var expense models.PersonalExpense
	err := r.Db.QueryRowContext(ctx, "SELECT id, amount, reason, description, category_id ,date FROM personal_expenses WHERE id = ?", id).
		Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			return expense, models.ErrExpenseNotFound
		}
		return expense, err
	}

	return expense, nil
}

// GetAllPersonalExpenses retrieves all expenses from the database.
func (r *PersonalExpenseRepository) GetAllPersonalExpenses(ctx context.Context) ([]models.PersonalExpense, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, reason, description, category_id, date FROM personal_expenses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PersonalExpense
	for rows.Next() {
		var expense models.PersonalExpense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *PersonalExpenseRepository) GetAllPersonalExpensesSummary(ctx context.Context) (*models.PersonalExpenseSummary, error) {
	// Query to get all expenses
	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, reason, description, category_id, date FROM personal_expenses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PersonalExpense
	var yearTotal, monthlyTotal float64

	// Get current year and month
	now := time.Now()
	currentYear, currentMonth := now.Year(), now.Month()

	for rows.Next() {
		var expense models.PersonalExpense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)

		// Parse the expense date
		expenseDate, err := time.Parse(time.RFC3339, expense.Date)
		if err != nil {
			return nil, err
		}

		// Check if the expense is in the current year
		if expenseDate.Year() == currentYear {
			yearTotal += expense.Amount

			// Check if the expense is also in the current month
			if expenseDate.Month() == currentMonth {
				monthlyTotal += expense.Amount
			}
		}
	}

	// Error check for rows iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Populate summary with totals and list of expenses
	summary := &models.PersonalExpenseSummary{
		MonthlyTotal: monthlyTotal,
		YearTotal:    yearTotal,
	}

	return summary, nil
}

func (r *PersonalExpenseRepository) GetPersonalExpensesSummaryBySubCategory(ctx context.Context, category_id int) (*models.PersonalExpenseSummary, error) {
	// Query to get all expenses with the specified category_id
	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, reason, description, category_id, date FROM personal_expenses WHERE category_id = ?", category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PersonalExpense
	var yearTotal, monthlyTotal float64

	// Get current year and month
	now := time.Now()
	currentYear, currentMonth := now.Year(), now.Month()

	for rows.Next() {
		var expense models.PersonalExpense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)

		// Parse the expense date
		expenseDate, err := time.Parse(time.RFC3339, expense.Date)
		if err != nil {
			return nil, err
		}

		// Check if the expense is in the current year
		if expenseDate.Year() == currentYear {
			yearTotal += expense.Amount

			// Check if the expense is also in the current month
			if expenseDate.Month() == currentMonth {
				monthlyTotal += expense.Amount
			}
		}
	}

	// Error check for rows iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Populate summary with totals
	summary := &models.PersonalExpenseSummary{
		MonthlyTotal: monthlyTotal,
		YearTotal:    yearTotal,
	}

	return summary, nil
}

func (r *PersonalExpenseRepository) GetPersonalExpensesSummaryByCategory(ctx context.Context, category_id int) (*models.PersonalExpenseSummary, error) {
	// Query to get all expenses with the specified category_id
	rows, err := r.Db.QueryContext(ctx, "SELECT personal_expenses.id, amount, reason, description, category_id, date FROM personal_expenses JOIN tender.categories ON personal_expenses.category_id = categories.id WHERE parent_id = ? OR category_id = ?", category_id, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PersonalExpense
	var yearTotal, monthlyTotal float64

	// Get current year and month
	now := time.Now()
	currentYear, currentMonth := now.Year(), now.Month()

	for rows.Next() {
		var expense models.PersonalExpense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
		if err != nil {
			return nil, err
		}

		expenses = append(expenses, expense)

		// Parse the expense date
		expenseDate, err := time.Parse(time.RFC3339, expense.Date)
		if err != nil {
			return nil, err
		}

		// Check if the expense is in the current year
		if expenseDate.Year() == currentYear {
			yearTotal += expense.Amount

			// Check if the expense is also in the current month
			if expenseDate.Month() == currentMonth {
				monthlyTotal += expense.Amount
			}
		}
	}

	// Error check for rows iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Populate summary with totals
	summary := &models.PersonalExpenseSummary{
		MonthlyTotal: monthlyTotal,
		YearTotal:    yearTotal,
	}

	return summary, nil
}

// GetAllPersonalExpenses retrieves all expenses from the database.
func (r *PersonalExpenseRepository) GetPersonalExpensesByCategoryId(ctx context.Context, category_id int) ([]models.PersonalExpense, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, amount, reason, description, category_id, date FROM personal_expenses WHERE category_id = ?", category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.PersonalExpense
	for rows.Next() {
		var expense models.PersonalExpense
		err := rows.Scan(&expense.ID, &expense.Amount, &expense.Reason, &expense.Description, &expense.CategoryID, &expense.Date)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

// UpdatePersonalExpense updates an existing expense in the database.
func (r *PersonalExpenseRepository) UpdatePersonalExpense(ctx context.Context, expense models.PersonalExpense) (models.PersonalExpense, error) {
	query := "UPDATE personal_expenses SET"
	params := []interface{}{}

	if expense.Amount != 0 {
		query += " amount = ?,"
		params = append(params, expense.Amount)
	}
	if expense.Reason != "" {
		query += " reason = ?,"
		params = append(params, expense.Reason)
	}
	if expense.Description != "" {
		query += " description = ?,"
		params = append(params, expense.Description)
	}
	if expense.CategoryID != 0 {
		query += " category_id = ?,"
		params = append(params, expense.CategoryID)
	}

	// Trim the last comma from the query
	query = query[:len(query)-1]
	query += " WHERE id = ?"
	params = append(params, expense.ID)

	_, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return models.PersonalExpense{}, err
	}

	// Retrieve the updated expense data
	row := r.Db.QueryRowContext(ctx, "SELECT id, amount, reason, description, category_id FROM personal_expenses WHERE id = ?", expense.ID)
	var updatedExpense models.PersonalExpense
	err = row.Scan(&updatedExpense.ID, &updatedExpense.Amount, &updatedExpense.Reason, &updatedExpense.Description, &updatedExpense.CategoryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.PersonalExpense{}, models.ErrExpenseNotFound
		}
		return models.PersonalExpense{}, err
	}

	return updatedExpense, nil
}

// DeletePersonalExpense removes an expense from the database by ID.
func (r *PersonalExpenseRepository) DeletePersonalExpense(ctx context.Context, id int) error {
	result, err := r.Db.ExecContext(ctx, "DELETE FROM personal_expenses WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return models.ErrExpenseNotFound
	}

	return nil
}
