package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type UserSurveyParticipation struct {
	gorm.Model
	UserId      uint      `gorm:"not null"`
	SurveyID    uint      `gorm:"not null"`
	StartAt     time.Time `gorm:"not null"`
	EndAt       time.Time
	CommittedAt time.Time
	User        models.User `gorm:"foreignKey:UserId;references:ID;"`
	Survey      Survey      `gorm:"foreignKey:SurveyID;references:ID;"`
}
