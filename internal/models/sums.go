package models

type Sums struct {
	TransactionsSum       float64 `json:"transactions_sum"`
	AdditionalExpensesSum float64 `json:"additional_expenses_sum"`
	TendersSum            float64 `json:"tenders_sum"`
}
