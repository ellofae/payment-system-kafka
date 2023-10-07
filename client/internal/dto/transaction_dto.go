package dto

import "time"

type TransactionData struct {
	UserID        int       `json:"user_id"`
	TransactionID string    `json:"transaction_id"`
	CardNumber    string    `json:"card_number" validate:"required,lte=20" form:"card_number" binding:"required"`
	Description   string    `json:"description" form:"description" binding:"required"`
	Amount        float64   `json:"amount" validate:"required,gt=0" form:"amount" binding:"required"`
	Date          time.Time `json:"date"`
}

type TransactionDisplayData struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	TransactionID string  `json:"transaction_id"`
	CardNumber    string  `json:"card_number"`
	Amount        float64 `json:"amount"`
	Date          string  `json:"date"`
}
