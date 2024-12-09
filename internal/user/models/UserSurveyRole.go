package models

import "time"

type UserSurveyRole struct {
	ID        uint
	UserID    uint `gorm:"not null"`
	SurveyID  uint `gorm:"not null"`
	RoleID    uint `gorm:"not null"`
	ExpiresAt time.Time
	User      User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Role      Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;"`
}
