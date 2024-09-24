package models

type Transaction struct {
	ID                int       `json:"id"`
	TransactionNumber int       `json:"transaction_number"`
	Type              string    `json:"type"`
	TenderNumber      *string   `json:"tender_number,omitempty"`
	UserID            int       `json:"user_id"`
	CompanyID         *int      `json:"company_id,omitempty"`
	Organization      *string   `json:"organization,omitempty"`
	Amount            float64   `json:"amount"`
	Total             float64   `json:"total"`
	Sell              float64   `json:"sell"`
	ProductName       string    `json:"product_name"`
	CompletedDate     string    `json:"completed_date"`
	Date              string    `json:"date"`
	Status            int       `json:"status"`
	Expenses          []Expense `json:"expenses"`
	UserName          *string   `json:"username,omitempty"`
	CompanyName       *string   `json:"companyname,omitempty"`
}

type Expense struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	TransactionID int     `json:"transaction_id"`
	Date          string  `json:"date"`
}
