package models

type Question struct {
	ID                uint
	SurveyID          uint   `gorm:"not null"`
	Text              string `gorm:"not null"`
	HasMultipleChoice bool   `gorm:"default:false"`
	MediaUrl          string
	Order             int
	LinkedQuestionID  uint
	Survey            Survey   `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;"`
	Choices           []Choice `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
}
