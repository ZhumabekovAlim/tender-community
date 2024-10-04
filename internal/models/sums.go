package models

type Sums struct {
	TransactionsSum       float64 `json:"transactions_sum"`
	AdditionalExpensesSum float64 `json:"additional_expenses_sum"`
	TendersGoikSum        float64 `json:"tenders_goik_sum"`
	TendersGoppSum        float64 `json:"tenders_gopp_sum"`
	ExtraTransactionsSum  float64 `json:"extra_transactions_sum"`
}
