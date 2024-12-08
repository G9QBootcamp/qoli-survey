package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	RoleID        uint           `gorm:"not null"`
	NationalID    string         `gorm:"unique;not null" json:"national_id"`
	Email         string         `gorm:"unique;not null" json:"email"`
	PasswordHash  string         `gorm:"not null"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	DateOfBirth   time.Time      `json:"date_of_birth"`
	City          string         `json:"city"`
	WalletBalance float64        `gorm:"default:0"`
	EmailVerified bool           `gorm:"default:false"`
	MaxSurveys    int
	GlobalRole    Role `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE;"`
}
