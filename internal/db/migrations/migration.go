package migrations

import (
	authModels "github.com/G9QBootcamp/qoli-survey/internal/auth/models"
	notificationModels "github.com/G9QBootcamp/qoli-survey/internal/notification/models"
	surveyModels "github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	userModels "github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&userModels.User{},
		&authModels.OTP{},
		&userModels.Role{},
		&userModels.Permission{},
		&surveyModels.Survey{},
		&userModels.UserSurveyRole{},
		&userModels.VoteVisibility{},
		&surveyModels.Question{},
		&surveyModels.Choice{},
		&surveyModels.Vote{},
		&surveyModels.UserSurveyParticipation{},
		&notificationModels.Notification{},
		&userModels.Transaction{},
	)
}
