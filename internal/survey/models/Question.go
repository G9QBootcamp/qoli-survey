package models

type Question struct {
	ID                uint   `json:"question_id"`
	SurveyID          uint   `gorm:"not null"`
	Text              string `gorm:"not null" json:"text"`
	HasMultipleChoice bool   `gorm:"default:false" json:"has_multiple_choice"`
	MediaUrl          string `json:"media_url"`
	Order             int
	LinkedQuestionID  uint
	Survey            Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;"`
	Choices           []Choice `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;" json:"choices"`
}
