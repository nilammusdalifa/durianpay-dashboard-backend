package entity

import "time"

type Payment struct {
	ID        string    `json:"id"`
	Merchant  string    `json:"merchant"`
	Status    string    `json:"status"`
	Amount    string    `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentFilter struct {
	Status *string
	ID     *string
	Sort   *string
}
