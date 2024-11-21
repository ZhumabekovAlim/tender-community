package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"tender/internal/models"
)

type ExtraTransactionRepository struct {
	Db *sql.DB
}

func (r *ExtraTransactionRepository) CreateExtraTransaction(ctx context.Context, extraTransaction models.ExtraTransaction) (models.ExtraTransaction, error) {
	result, err := r.Db.ExecContext(ctx, `
		INSERT INTO extra_transactions (user_id, description, total, status)
		VALUES (?, ?, ?, ?)`,
		extraTransaction.UserID, extraTransaction.Description, extraTransaction.Total, extraTransaction.Status)
	if err != nil {
		return models.ExtraTransaction{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.ExtraTransaction{}, err
	}
	extraTransaction.ID = int(id)

	return extraTransaction, nil
}

func (r *ExtraTransactionRepository) GetExtraTransactionByID(ctx context.Context, id int) (models.ExtraTransaction, error) {
	var extraTransaction models.ExtraTransaction
	err := r.Db.QueryRowContext(ctx, `
		SELECT et.id, user_id, description, total, date, status,CONCAT(u.name, ' ', u.last_name) as username
		FROM extra_transactions et
		JOIN tender.users u on u.id = et.user_id
		WHERE et.id = ?`, id).
		Scan(&extraTransaction.ID, &extraTransaction.UserID, &extraTransaction.Description, &extraTransaction.Total, &extraTransaction.Date, &extraTransaction.Status, &extraTransaction.UserName)
	if err != nil {
		return models.ExtraTransaction{}, err
	}

	return extraTransaction, nil
}

func (r *ExtraTransactionRepository) GetAllExtraTransactions(ctx context.Context) ([]models.ExtraTransaction, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT et.id, user_id, description, total, date, status, u.name 
		FROM extra_transactions et
		JOIN tender.users u on u.id = et.user_id
		ORDER BY date DESC`)
	if err != nil {
		log.Printf("Error querying extra transactions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var extraTransaction models.ExtraTransaction
		err := rows.Scan(&extraTransaction.ID, &extraTransaction.UserID, &extraTransaction.Description, &extraTransaction.Total, &extraTransaction.Date, &extraTransaction.Status, &extraTransaction.UserName)
		if err != nil {
			log.Printf("Error scanning extra transaction row: %v", err)
			return nil, err
		}
		extraTransactions = append(extraTransactions, extraTransaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over extra transactions rows: %v", err)
		return nil, err
	}

	return extraTransactions, nil
}

func (r *ExtraTransactionRepository) GetExtraTransactionsByUser(ctx context.Context, userID int) ([]models.ExtraTransaction, error) {
	rows, err := r.Db.QueryContext(ctx, `
		SELECT extra_transactions.id, user_id, description, total, date, extra_transactions.status , CONCAT(u.name, ' ', u.last_name) as username
		FROM extra_transactions
		JOIN tender.users u on extra_transactions.user_id = u.id
		WHERE user_id = ? ORDER BY date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var extraTransaction models.ExtraTransaction
		err := rows.Scan(&extraTransaction.ID, &extraTransaction.UserID, &extraTransaction.Description, &extraTransaction.Total, &extraTransaction.Date, &extraTransaction.Status, &extraTransaction.UserName)
		if err != nil {
			return nil, err
		}
		extraTransactions = append(extraTransactions, extraTransaction)
	}

	return extraTransactions, nil
}

func (r *ExtraTransactionRepository) UpdateExtraTransaction(ctx context.Context, extraTransaction models.ExtraTransaction) (models.ExtraTransaction, error) {
	query := "UPDATE extra_transactions SET"
	params := []interface{}{}

	if extraTransaction.UserID != 0 {
		query += " user_id = ?,"
		params = append(params, extraTransaction.UserID)
	}
	if extraTransaction.Description != "" {
		query += " description = ?,"
		params = append(params, extraTransaction.Description)
	}
	if extraTransaction.Total != 0 {
		query += " total = ?,"
		params = append(params, extraTransaction.Total)
	}
	if extraTransaction.Status != 0 {
		query += " status = ?,"
		params = append(params, extraTransaction.Status)
	}

	// Remove the trailing comma and add the WHERE clause
	query = query[:len(query)-1] + " WHERE id = ?"
	params = append(params, extraTransaction.ID)

	_, err := r.Db.ExecContext(ctx, query, params...)
	if err != nil {
		return models.ExtraTransaction{}, err
	}

	// Retrieve the updated record from the database
	var updatedTransaction models.ExtraTransaction
	err = r.Db.QueryRowContext(ctx, `
		SELECT id, user_id, description, total, status, date
		FROM extra_transactions
		WHERE id = ?`, extraTransaction.ID).
		Scan(&updatedTransaction.ID, &updatedTransaction.UserID, &updatedTransaction.Description,
			&updatedTransaction.Total, &updatedTransaction.Status, &updatedTransaction.Date)

	if err != nil {
		return models.ExtraTransaction{}, err
	}

	return updatedTransaction, nil
}

func (r *ExtraTransactionRepository) DeleteExtraTransaction(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM extra_transactions WHERE id = ?`, id)
	return err
}

func (r *ExtraTransactionRepository) GetExtraTransactionCountsByUserID(ctx context.Context, userID int) (*models.ExtraTransactionCount, error) {
	// Query to count total extra transactions by user ID
	queryTotal := `
        SELECT COUNT(*) 
        FROM extra_transactions 
        WHERE user_id = ?;
    `

	// Query to count extra transactions by user ID and status
	queryStatus := `
        SELECT COUNT(*) 
        FROM extra_transactions 
        WHERE user_id = ? AND status = ?;
    `

	counts := &models.ExtraTransactionCount{}

	// Execute total transactions query
	err := r.Db.QueryRowContext(ctx, queryTotal, userID).Scan(&counts.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to count total extra transactions: %w", err)
	}

	// Execute queries for each status
	statusCounts := []*int{&counts.Status0, &counts.Status1, &counts.Status2, &counts.Status3}
	for i := 0; i < 4; i++ {
		err = r.Db.QueryRowContext(ctx, queryStatus, userID, i).Scan(statusCounts[i])
		if err != nil {
			return nil, fmt.Errorf("failed to count status %d extra transactions: %w", i, err)
		}
	}

	return counts, nil
}

func (r *ExtraTransactionRepository) GetAllExtraTransactionsByDateRange(ctx context.Context, startDate, endDate string, userId int) ([]models.ExtraTransaction, error) {
	var query string
	var rows *sql.Rows
	var err error

	if userId == 1 {
		query = `
			SELECT e.id, e.user_id, e.description, e.total, e.date, e.status, u.name as username
			FROM extra_transactions e 
			JOIN users u ON u.id = e.user_id
			WHERE e.date BETWEEN ? AND ?
			ORDER BY e.date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, startDate, endDate)
	} else {
		query = `
			SELECT e.id, e.user_id, e.description, e.total, e.date, e.status, u.name as username
			FROM extra_transactions e 
			JOIN users u ON u.id = e.user_id
			WHERE e.user_id = ? AND e.date BETWEEN ? AND ?
			ORDER BY e.date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, userId, startDate, endDate)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query extra transactions: %w", err)
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var extra models.ExtraTransaction
		if err := rows.Scan(&extra.ID, &extra.UserID, &extra.Description, &extra.Total, &extra.Date,
			&extra.Status, &extra.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan extra transaction: %w", err)
		}
		extraTransactions = append(extraTransactions, extra)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return extraTransactions, nil
}

func (r *ExtraTransactionRepository) GetAllExtraTransactionsByDateRangeCompany(ctx context.Context, startDate, endDate string, userId, companyId int) ([]models.ExtraTransaction, error) {
	var query string
	var rows *sql.Rows
	var err error

	if userId == 1 {
		query = `
			SELECT e.id, e.user_id, e.description, e.total, e.date, e.status, u.name as username
			FROM extra_transactions e 
			JOIN users u ON u.id = e.user_id
			WHERE e.user_id = ? AND e.date BETWEEN ? AND ?
			ORDER BY e.date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, startDate, endDate, companyId)
	} else {
		query = `
			SELECT e.id, e.user_id, e.description, e.total, e.date, e.status, u.name as username
			FROM extra_transactions e 
			JOIN users u ON u.id = e.user_id
			WHERE e.user_id = ? AND e.user_id = ? AND e.date BETWEEN ? AND ?
			ORDER BY e.date DESC
		`
		rows, err = r.Db.QueryContext(ctx, query, userId, companyId, startDate, endDate)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query extra transactions: %w", err)
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var extra models.ExtraTransaction
		if err := rows.Scan(&extra.ID, &extra.UserID, &extra.Description, &extra.Total, &extra.Date,
			&extra.Status, &extra.UserName); err != nil {
			return nil, fmt.Errorf("failed to scan extra transaction: %w", err)
		}
		extraTransactions = append(extraTransactions, extra)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return extraTransactions, nil
}
