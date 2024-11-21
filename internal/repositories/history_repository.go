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
func (r *HistoryRepository) GetAllHistory(ctx context.Context, source *string, startDate, endDate string, limit, offset int) ([]models.CombinedAction, error) {
	query := `
		WITH combined_data AS (
			SELECT id, 'Transaction' AS source, amount, total, date
			FROM transactions
			UNION ALL
			SELECT id, 'Tender' AS source, total AS amount, NULL AS total, date
			FROM tenders
			UNION ALL
			SELECT id, 'PersonalExpense' AS source, amount, NULL AS total, date
			FROM personal_expenses
			UNION ALL
			SELECT id, 'PersonalDebt' AS source, amount, NULL AS total, get_date AS date
			FROM personal_debts
			UNION ALL
			SELECT id, 'ExtraTransaction' AS source, total AS amount, NULL AS total, date
			FROM extra_transactions
			UNION ALL
			SELECT id, 'BalanceHistory' AS source, amount, NULL AS total, created_at AS date
			FROM balance_history
		)
		SELECT *
		FROM combined_data
		WHERE
			date >= ?
			AND date <= ?
			AND (? IS NULL OR source = ?)
		ORDER BY date DESC
		LIMIT ? OFFSET ?;
	`

	rows, err := r.Db.QueryContext(ctx, query, startDate, endDate, source, source, limit, offset)
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
