package models

import "time"

type Transaction struct {
	ID         string    `json:"id"`
	UserID     string    `validate:"required" json:"user_id"`
	Amount     float64   `validate:"required" json:"amount"`
	CategoryID string    `validate:"required" json:"category_id"`
	Comment    string    `validate:"required" json:"comment"`
	Date       time.Time `validate:"required" json:"date"`
	CreatedAt  time.Time `json:"created_at"`
}
