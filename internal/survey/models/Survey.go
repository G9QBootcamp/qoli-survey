package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type Survey struct {
	ID                 uint `gorm:"primarykey" json:"survey_id"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt          `gorm:"index"`
	OwnerID            uint                    `gorm:"not null" json:"user_id"`
	Title              string                  `gorm:"not null" json:"title"`
	StartTime          time.Time               `gorm:"not null" json:"start_time"`
	EndTime            time.Time               `gorm:"not null" json:"end_time"`
	IsSequential       bool                    `gorm:"default:false" json:"is_sequential"`
	AllowReturn        bool                    `gorm:"default:false" json:"allow_return"`
	ParticipationLimit int                     `gorm:"default:1" json:"participation_limit"`
	AnswerTimeLimit    int                     `gorm:"not null" json:"answer_time_limit"`
	Owner              models.User             `gorm:"foreignKey:OwnerID;references:ID;constraint:OnDelete:CASCADE;"`
	Questions          []Question              `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;" json:"questions"`
	UserSurveyRoles    []models.UserSurveyRole `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;"`
	VoteVisibilities   []models.VoteVisibility `gorm:"foreignKey:SurveyID;constraint:OnDelete:CASCADE;"`
	Options            []SurveyOption          `gorm:"foreignKey:SurveyId;constraint:OnDelete:CASCADE;"`
}
