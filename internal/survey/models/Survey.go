package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type Survey struct {
	gorm.Model
	OwnerID            uint                    `gorm:"not null"`
	Title              string                  `gorm:"not null"`
	StartTime          time.Time               `gorm:"not null"`
	EndTime            time.Time               `gorm:"not null"`
	IsSequential       bool                    `gorm:"default:false"`
	AllowReturn        bool                    `gorm:"default:false"`
	ParticipationLimit int                     `gorm:"default:1"`
	AnswerTimeLimit    int                     `gorm:"not null"`
	Owner              models.User             `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`
	Questions          []Question              `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;"`
	UserSurveyRoles    []models.UserSurveyRole `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;"`
	VoteVisibilities   []models.VoteVisibility `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;"`
}
