package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"

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
	WalletBalance float64         `gorm:"default:0"`
	GlobalRole    Role            `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE;"`
	Surveys       []models.Survey `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE;"`
	Votes         []models.Vote   `gorm:"foreignKey:VoterID;constraint:OnDelete:CASCADE;"`
}
