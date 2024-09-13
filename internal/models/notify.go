package models

type Notify struct {
	ID       int    `json:"id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Sender   int    `json:"sender"`
	Receiver int    `json:"receiver"`
}
