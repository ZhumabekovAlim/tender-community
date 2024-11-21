package models

type CombinedAction struct {
	ID     int      `json:"id"`
	Source string   `json:"source"` // Table name or source
	Amount float64  `json:"amount"`
	Total  *float64 `json:"total,omitempty"` // Nullable
	Date   string   `json:"date"`
}

type HistoryRequest struct {
	Source    *string `json:"source,omitempty"` // Filter by source (optional)
	StartDate string  `json:"start_date"`       // Start date (required)
	EndDate   string  `json:"end_date"`         // End date (required)
	Limit     int     `json:"limit,omitempty"`  // Number of records to return (default: 10)
	Offset    int     `json:"offset,omitempty"` // Number of records to skip (default: 0)
}
