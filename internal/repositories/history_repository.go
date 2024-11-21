package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type HistoryRepository struct {
	Db *sql.DB
}

// GetAllHistory retrieves combined history data from all tables.
func (r *HistoryRepository) GetAllHistory(ctx context.Context) ([]models.CombinedAction, error) {
	query := `
	(
		SELECT id, 'Transaction' as source, amount, total, date
		FROM transactions
	)
	UNION ALL
	(
		SELECT id, 'Tender' as source, total as amount, NULL as total, date
		FROM tenders
	)
	UNION ALL
	(
		SELECT id, 'PersonalExpense' as source, amount, NULL as total, date
		FROM personal_expenses
	)
	UNION ALL
	(
		SELECT id, 'PersonalDebt' as source, amount, NULL as total, get_date as date
		FROM personal_debts
	)
	UNION ALL
	(
		SELECT id, 'ExtraTransaction' as source, total as amount, NULL as total, date
		FROM extra_transactions
	)
	UNION ALL
	(
		SELECT id, 'BalanceHistory' as source, amount, NULL as total, created_at as date
		FROM balance_history
	)
	ORDER BY date DESC
	LIMIT 50;
	`

	rows, err := r.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.CombinedAction
	for rows.Next() {
		var action models.CombinedAction
		var total sql.NullFloat64 // Handle NULL for total
		err := rows.Scan(&action.ID, &action.Source, &action.Amount, &total, &action.Date)
		if err != nil {
			return nil, err
		}

		if total.Valid {
			action.Total = &total.Float64
		} else {
			action.Total = nil
		}

		history = append(history, action)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return history, nil
}
