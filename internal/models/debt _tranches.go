package models

type DebtTranche struct {
	ID          int     `json:"id"`
	DebtID      int     `json:"debt_id"`
	Amount      float64 `json:"amount"`
	Description *string `json:"description,omitempty"`
	Date        string  `json:"date"`
}
