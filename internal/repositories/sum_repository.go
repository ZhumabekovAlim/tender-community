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

func (r *SumRepository) GetDebtsByAccount(ctx context.Context) ([]models.AccountDebts, error) {
	query := `
    WITH tender_numbers AS (
        SELECT tender_number FROM transactions WHERE status = 2
        UNION
        SELECT t.tender_number FROM additional_expenses ae
        JOIN transactions t ON ae.transaction_id = t.id
        WHERE t.status = 2
        UNION
        SELECT tender_number FROM tenders WHERE status = 2
    )
    SELECT
        tn.tender_number AS account_number,
        COALESCE(ts.transactions_sum, 0) AS transactions_sum,
        COALESCE(aes.additional_expenses_sum, 0) AS additional_expenses_sum,
        COALESCE(goik.tenders_goik_sum, 0) AS tenders_goik_sum,
        COALESCE(gopp.tenders_gopp_sum, 0) AS tenders_gopp_sum
    FROM tender_numbers tn
    LEFT JOIN (
        SELECT tender_number, SUM(total) AS transactions_sum
        FROM transactions
        WHERE status = 2
        GROUP BY tender_number
    ) ts ON tn.tender_number = ts.tender_number
    LEFT JOIN (
        SELECT t.tender_number, SUM(ae.amount) AS additional_expenses_sum
        FROM additional_expenses ae
        JOIN transactions t ON ae.transaction_id = t.id
        WHERE t.status = 2
        GROUP BY t.tender_number
    ) aes ON tn.tender_number = aes.tender_number
    LEFT JOIN (
        SELECT tender_number, SUM(total) AS tenders_goik_sum
        FROM tenders
        WHERE status = 2 AND type = 'ГОИК'
        GROUP BY tender_number
    ) goik ON tn.tender_number = goik.tender_number
    LEFT JOIN (
        SELECT tender_number, SUM(total) AS tenders_gopp_sum
        FROM tenders
        WHERE status = 2 AND type = 'ГОПП'
        GROUP BY tender_number
    ) gopp ON tn.tender_number = gopp.tender_number
    ORDER BY tn.tender_number
    `

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var debts []models.AccountDebts
	for rows.Next() {
		var debt models.AccountDebts
		err := rows.Scan(
			&debt.AccountNumber,
			&debt.TransactionsSum,
			&debt.AdditionalExpensesSum,
			&debt.TendersGoikSum,
			&debt.TendersGoppSum,
		)
		if err != nil {
			return nil, err
		}
		debts = append(debts, debt)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return debts, nil
}
