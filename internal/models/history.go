package models

type CombinedAction struct {
	ID     int      `json:"id"`
	Source string   `json:"source"` // Table name or source
	Amount float64  `json:"amount"`
	Total  *float64 `json:"total,omitempty"` // Nullable
	Date   string   `json:"date"`
}
