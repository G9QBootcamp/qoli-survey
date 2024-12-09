package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID uint    `gorm:"not null"`
	Amount float64 `gorm:"not null"`
	Reason string  `gorm:"not null"`
	User   User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
