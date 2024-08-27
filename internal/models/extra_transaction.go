package models

type ExtraTransaction struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Description string  `json:"description"`
	Total       float64 `json:"total"`
	Date        string  `json:"date"`
	Status      int     `json:"status"`
	UserName    string  `json:"name"`
}
