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

type ExtraTransactionCount struct {
	Total   int `json:"total"`
	Status0 int `json:"status_0"`
	Status1 int `json:"status_1"`
	Status2 int `json:"status_2"`
	Status3 int `json:"status_3"`
}
