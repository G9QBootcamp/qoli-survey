package models

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

type Vote struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	VoterID    uint           `gorm:"not null" json:"voter_id"`
	QuestionID uint           `gorm:"not null" json:"question_id"`
	Answer     string         `gorm:"not null" json:"answer"`
	IsCorrect  bool           `json:"is_correct"`
	Voter      models.User    `gorm:"foreignKey:VoterID;references:ID;constraint:OnDelete:CASCADE;" json:"voter"`
	Question   Question       `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;" json:"question"`
}
