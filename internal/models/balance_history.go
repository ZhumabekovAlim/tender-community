package models

type BalanceHistory struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	UserID      int     `json:"user_id"`
	CategoryID  int     `json:"category_id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
