package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	BuyerID  uint    `gorm:"not null"`
	SellerID uint    `gorm:"not null"`
	Amount   float64 `gorm:"not null"`
}
