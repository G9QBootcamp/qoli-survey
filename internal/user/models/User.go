package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	NationalID    string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	PasswordHash  string `gorm:"not null"`
	FirstName     string
	LastName      string
	DateOfBirth   time.Time
	City          string
	WalletBalance float64 `gorm:"default:0"`
}
