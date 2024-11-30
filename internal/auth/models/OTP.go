package models

import (
	"time"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	ID        uint
	UserID    uint      `gorm:"not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsValid   bool      `gorm:"default:true"`
}
