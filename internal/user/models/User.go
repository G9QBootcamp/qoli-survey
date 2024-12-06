package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	RoleID        uint   `gorm:"not null"`
	NationalID    string `gorm:"unique;not null"`
	Email         string `gorm:"unique;not null"`
	PasswordHash  string `gorm:"not null"`
	FirstName     string
	LastName      string
	DateOfBirth   time.Time
	City          string
	EmailVerified bool    `gorm:"default:false"`
	WalletBalance float64 `gorm:"default:0"`
	GlobalRole    Role    `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE;"`
}
