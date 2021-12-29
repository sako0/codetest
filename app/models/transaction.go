package models

type Transaction struct {
	ID          int    `json:"id"`
	UserId      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}
