package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type UserSurveyParticipation struct {
	ID          uint `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	UserId      uint           `gorm:"not null" json:"user_id"`
	SurveyID    uint           `gorm:"not null" json:"survey_id"`
	StartAt     time.Time      `gorm:"not null" json:"start_at"`
	EndAt       *time.Time     `gorm:"default:null" json:"end_at"`
	CommittedAt *time.Time     `gorm:"default:null" json:"committed_at"`
	User        models.User    `gorm:"foreignKey:UserId;references:ID;"`
	Survey      Survey         `gorm:"foreignKey:SurveyID;references:ID;"`
}
