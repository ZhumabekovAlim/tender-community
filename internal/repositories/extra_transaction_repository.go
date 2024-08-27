package repositories

import (
	"context"
	"database/sql"
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
		SELECT id, user_id, description, total, date, status 
		FROM extra_transactions WHERE id = ?`, id).
		Scan(&extraTransaction.ID, &extraTransaction.UserID, &extraTransaction.Description, &extraTransaction.Total, &extraTransaction.Date, &extraTransaction.Status)
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
		SELECT id, user_id, description, total, date, status 
		FROM extra_transactions WHERE user_id = ? ORDER BY date DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var extraTransactions []models.ExtraTransaction
	for rows.Next() {
		var extraTransaction models.ExtraTransaction
		err := rows.Scan(&extraTransaction.ID, &extraTransaction.UserID, &extraTransaction.Description, &extraTransaction.Total, &extraTransaction.Date, &extraTransaction.Status)
		if err != nil {
			return nil, err
		}
		extraTransactions = append(extraTransactions, extraTransaction)
	}

	return extraTransactions, nil
}

func (r *ExtraTransactionRepository) UpdateExtraTransaction(ctx context.Context, extraTransaction models.ExtraTransaction) (models.ExtraTransaction, error) {
	_, err := r.Db.ExecContext(ctx, `
		UPDATE extra_transactions SET user_id = ?, description = ?, total = ?, status = ? 
		WHERE id = ?`,
		extraTransaction.UserID, extraTransaction.Description, extraTransaction.Total, extraTransaction.Status, extraTransaction.ID)
	if err != nil {
		return models.ExtraTransaction{}, err
	}

	return extraTransaction, nil
}

func (r *ExtraTransactionRepository) DeleteExtraTransaction(ctx context.Context, id int) error {
	_, err := r.Db.ExecContext(ctx, `DELETE FROM extra_transactions WHERE id = ?`, id)
	return err
}
