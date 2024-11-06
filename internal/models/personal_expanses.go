package models

type PersonalExpense struct {
	ID          int     `json:"id"`
	Amount      float64 `json:"amount"`
	Reason      string  `json:"reason"`
	Description string  `json:"description,omitempty"`
	CategoryID  int     `json:"category_id"`
	Date        string  `json:"date"`
}

type PersonalExpenseSummary struct {
	Expenses     []PersonalExpense `json:"expenses"`
	MonthlyTotal float64           `json:"monthly_total"`
	AllTimeTotal float64           `json:"all_time_total"`
}
