package models

type AccountDebts struct {
	AccountNumber         string  `json:"account_number"`
	TransactionsSum       float64 `json:"transactions_sum"`
	AdditionalExpensesSum float64 `json:"additional_expenses_sum"`
	TendersGoikSum        float64 `json:"tenders_goik_sum"`
	TendersGoppSum        float64 `json:"tenders_gopp_sum"`
}
