package models

import "time"

// TransactionData represents data for a single transaction
type TransactionData struct {
	TransactionNumber *string `json:"transaction_number,omitempty"`
	Amount            float64 `json:"amount"`
}

// TenderData represents data for a single tender
type TenderData struct {
	TenderNumber *string `json:"tender_number"`
	Amount       float64 `json:"amount"`
}

// AdditionalExpenseData represents data for a single additional expense
type AdditionalExpenseData struct {
	Date   time.Time `json:"date,omitempty"`
	Amount float64   `json:"amount"`
}

// ClientData aggregates all the data for a client
type ClientData struct {
	UserID             int                     `json:"user_id"`
	Transactions       []TransactionData       `json:"transactions"`
	TendersGOIK        []TenderData            `json:"tenders_goik"`
	TendersGOPP        []TenderData            `json:"tenders_gopp"`
	AdditionalExpenses []AdditionalExpenseData `json:"additional_expenses"`
}
