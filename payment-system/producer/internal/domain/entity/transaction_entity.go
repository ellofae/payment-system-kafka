package entity

type TransactionData struct {
	UserID        int     `json:"user_id"`
	TransactionID string  `json:"transaction_id"`
	CardNumber    string  `json:"card_number"`
	Description   string  `json:"description"`
	Amount        float64 `json:"amount"`
}
