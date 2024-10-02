package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type SumRepository struct {
	Db *sql.DB
}

// GetSumsByUserID retrieves the sums from all tables for a given user_id and status = 2.
func (r *SumRepository) GetSumsByUserID(ctx context.Context, userID int) (models.Sums, error) {
	query := `
    SELECT
        (SELECT COALESCE(SUM(total), 0) FROM transactions WHERE user_id = ? AND status = 2) AS transactions_sum,
        (SELECT COALESCE(SUM(ae.amount), 0)
         FROM additional_expenses ae
         JOIN transactions t ON ae.transaction_id = t.id
         WHERE t.user_id = ? AND t.status = 2) AS additional_expenses_sum,
        (SELECT COALESCE(SUM(total), 0) FROM tenders WHERE user_id = ? AND status = 2 AND type = 'ГОИК') AS tenders_goik_sum,
        (SELECT COALESCE(SUM(total), 0) FROM tenders WHERE user_id = ? AND status = 2 AND type = 'ГОПП') AS tenders_gopp_sum
    `
	row := r.Db.QueryRowContext(ctx, query, userID, userID, userID, userID)
	var sums models.Sums
	err := row.Scan(&sums.TransactionsSum, &sums.AdditionalExpensesSum, &sums.TendersGoikSum, &sums.TendersGoppSum)
	if err != nil {
		return sums, err
	}
	return sums, nil
}
