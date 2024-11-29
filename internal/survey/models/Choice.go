package models

type Choice struct {
	ID               uint
	QuestionID       uint   `gorm:"not null"`
	Text             string `gorm:"unique;not null"`
	IsCorrect        bool
	LinkedQuestionID uint
	Question         Question `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"`
}
