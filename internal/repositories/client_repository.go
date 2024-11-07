package repositories

import (
	"context"
	"database/sql"
	"tender/internal/models"
)

type ClientRepository struct {
	Db *sql.DB
}

// GetClientData retrieves data for a specific client (user_id)
func (r *ClientRepository) GetClientData(ctx context.Context, userID int) (models.ClientData, error) {
	var clientData models.ClientData
	clientData.UserID = userID

	// Transactions data
	transactionsQuery := `
    SELECT transaction_number, sell
    FROM transactions
    WHERE user_id = ? AND status = 2
    `
	transRows, err := r.Db.QueryContext(ctx, transactionsQuery, userID)
	if err != nil {
		return clientData, err
	}
	defer transRows.Close()

	for transRows.Next() {
		var transaction models.TransactionData
		err := transRows.Scan(&transaction.TransactionNumber, &transaction.Amount)
		if err != nil {
			return clientData, err
		}
		var empty *string = new(string)
		*empty = "-"
		if transaction.TransactionNumber == nil {
			transaction.TransactionNumber = empty
		}
		clientData.Transactions = append(clientData.Transactions, transaction)
	}
	if err := transRows.Err(); err != nil {
		return clientData, err
	}

	// Tenders GOIK data
	tendersGOIKQuery := `
    SELECT tender_number, total
    FROM tenders
    WHERE user_id = ? AND status = 2 AND type = 'ГОИК'
    `
	goikRows, err := r.Db.QueryContext(ctx, tendersGOIKQuery, userID)
	if err != nil {
		return clientData, err
	}
	defer goikRows.Close()

	for goikRows.Next() {
		var tender models.TenderData
		err := goikRows.Scan(&tender.TenderNumber, &tender.Amount)
		if err != nil {
			return clientData, err
		}
		clientData.TendersGOIK = append(clientData.TendersGOIK, tender)
	}
	if err := goikRows.Err(); err != nil {
		return clientData, err
	}

	// Tenders GOPP data
	tendersGOPPQuery := `
    SELECT tender_number, total
    FROM tenders
    WHERE user_id = ? AND status = 2 AND type = 'ГОПП'
    `
	goppRows, err := r.Db.QueryContext(ctx, tendersGOPPQuery, userID)
	if err != nil {
		return clientData, err
	}
	defer goppRows.Close()

	for goppRows.Next() {
		var tender models.TenderData
		err := goppRows.Scan(&tender.TenderNumber, &tender.Amount)
		if err != nil {
			return clientData, err
		}
		clientData.TendersGOPP = append(clientData.TendersGOPP, tender)
	}
	if err := goppRows.Err(); err != nil {
		return clientData, err
	}

	// Additional expenses data
	expensesQuery := `
    SELECT ae.date, ae.amount
    FROM additional_expenses ae
    JOIN transactions t ON ae.transaction_id = t.id
    WHERE t.user_id = ? AND t.status = 2
    `
	expenseRows, err := r.Db.QueryContext(ctx, expensesQuery, userID)
	if err != nil {
		return clientData, err
	}
	defer expenseRows.Close()

	for expenseRows.Next() {
		var expense models.AdditionalExpenseData
		err := expenseRows.Scan(&expense.Date, &expense.Amount)
		if err != nil {
			return clientData, err
		}
		clientData.AdditionalExpenses = append(clientData.AdditionalExpenses, expense)
	}
	if err := expenseRows.Err(); err != nil {
		return clientData, err
	}
	//a := clientData.Transactions. + clientData.AdditionalExpenses

	return clientData, nil
}
