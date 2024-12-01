package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"

	"gorm.io/gorm"
)

type OTP struct {
	gorm.Model
	UserID    uint        `gorm:"not null"`
	Code      string      `gorm:"not null"`
	ExpiresAt time.Time   `gorm:"not null"`
	IsValid   bool        `gorm:"default:true"`
	User      models.User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
}
