package models

type PersonalDebt struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	Type       string  `json:"type"`
	GetDate    *string `json:"get_date,omitempty"`
	ReturnDate *string `json:"return_date,omitempty"`
	Status     int     `json:"status"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
