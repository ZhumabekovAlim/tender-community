package models

type BalanceCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ParentID  int    `json:"parent_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
