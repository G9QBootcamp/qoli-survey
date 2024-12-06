package models

type Choice struct {
	ID               uint     `json:"choice_id"`
	QuestionID       uint     `gorm:"not null"`
	Text             string   `gorm:"not null" json:"text"`
	IsCorrect        bool     `json:"is_correct"`
	LinkedQuestionID uint     `json:"linked_question_id"`
	Question         Question `gorm:"foreignKey:QuestionID;references:ID;constraint:OnDelete:CASCADE;"`
}
