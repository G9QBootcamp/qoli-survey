package models

import "github.com/G9QBootcamp/qoli-survey/internal/survey/models"

type UserSurveyRole struct {
	ID       uint
	UserID   uint `gorm:"not null"`
	SurveyID uint `gorm:"not null"`
	RoleID   uint `gorm:"not null"`
	duration int
	User     User          `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;"`
	Survey   models.Survey `gorm:"foreignKey:SurveyID;references:ID;constraint:OnDelete:CASCADE;"`
	Role     Role          `gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE;"`
}
