package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type SurveyOption struct {
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserId    uint           `gorm:"not null" json:"user_id"`
	SurveyId  uint           `gorm:"not null" json:"survey_id"`
	Name      string         `gorm:"not null" json:"name"`
	Value     string         `gorm:"not null" json:"value"`
	User      models.User    `gorm:"foreignKey:UserId;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	Survey    Survey         `gorm:"foreignKey:SurveyId;references:ID;constraint:OnDelete:CASCADE;" json:"survey"`
}
