package models

import "time"

type Transaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BuyerID   uint      `json:"buyer_id"`
	SellerID  uint      `json:"seller_id"`
	VoteCount int       `json:"vote_count"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
