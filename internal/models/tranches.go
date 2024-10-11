package models

type Tranche struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	Amount        int    `json:"amount"`
	Description   string `json:"description"`
	Date          string `json:"date"`
}
