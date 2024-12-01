package models

type UserSurveyRole struct {
	ID        uint
	UserID    uint `gorm:"not null"`
	SurveyID  uint `gorm:"not null"`
	RoleID    uint `gorm:"not null"`
	TimeLimit int
	User      User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Role      Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE;"`
}
