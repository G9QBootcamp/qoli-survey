package models

type VoteVisibility struct {
	ID           uint
	SurveyID     uint `gorm:"not null"`
	ViewerID     uint `gorm:"not null"` // User allowed to view votes
	RespondentID uint `gorm:"not null"` // User whose votes can be viewed
	Viewer       User `gorm:"foreignKey:ViewerID;constraint:OnDelete:CASCADE;"`
	Respondent   User `gorm:"foreignKey:RespondentID;constraint:OnDelete:CASCADE;"`
}
