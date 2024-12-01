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
	Duration   int         `gorm:"default:0"`
	Voter      models.User `gorm:"foreignKey:VoterID;references:ID;constraint:OnDelete:CASCADE;"`
	Question   Survey      `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"`
}
