package models

type Part struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
