package models

import (
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type Vote struct {
	gorm.Model
	VoterID    uint        `gorm:"not null"`
	QuestionID uint        `gorm:"not null"`
	Answer     string      `gorm:"not null"`
	IsCorrect  bool        `json:"is_correct"`
	Voter      models.User `gorm:"foreignKey:VoterID;references:ID;constraint:OnDelete:CASCADE;"`
	Question   Question    `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"`
}
