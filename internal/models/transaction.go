package models

type Transaction struct {
	ID                int       `json:"id"`
	TransactionNumber *string   `json:"transaction_number,omitempty"`
	Type              string    `json:"type"`
	TenderNumber      *string   `json:"tender_number,omitempty"`
	UserID            *int      `json:"user_id,omitempty"`
	CompanyID         *int      `json:"company_id,omitempty"`
	Organization      *string   `json:"organization,omitempty"`
	Amount            float64   `json:"amount"`
	Total             float64   `json:"total"`
	Sell              float64   `json:"sell"`
	ProductName       string    `json:"product_name"`
	CompletedDate     *string   `json:"completed_date,omitempty"`
	Date              string    `json:"date"`
	Status            int       `json:"status"`
	Expenses          []Expense `json:"expenses"`
	UserName          string    `json:"username,omitempty"`
	CompanyName       *string   `json:"companyname,omitempty"`
	Debt              float64   `json:"debt,omitempty"`
	Margin            *float64  `json:"margin,omitempty"`
}

type Expense struct {
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	TransactionID int     `json:"transaction_id"`
	Date          string  `json:"date"`
}

type TransactionDebt struct {
	Zakup float64 `json:"zakup"`
}

type TransactionDebtId struct {
	Debt float64 `json:"debt"`
}

type TransactionCount struct {
	TotalTransactions int
	Status0           int
	Status1           int
	Status2           int
	Status3           int
}
