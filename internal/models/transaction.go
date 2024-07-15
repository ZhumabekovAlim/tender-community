package models

import "time"

type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	CompanyID int       `json:"company_id"`
	Amount    float64   `json:"amount"`
	Date      time.Time `json:"date"`
	Status    int       `json:"status"`
}
