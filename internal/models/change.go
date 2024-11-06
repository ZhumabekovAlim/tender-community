package models

type Change struct {
	ID            int     `json:"id"`
	TransactionID int     `json:"transaction_id"`
	Amount        int     `json:"amount"`
	Description   *string `json:"description,omitempty"`
	Date          string  `json:"date"`
}
