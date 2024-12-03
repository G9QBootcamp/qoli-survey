package seeds

import (
	"log"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type Seeder struct {
	db     db.DbService
	logger logging.Logger
}

func NewSeeder(db db.DbService, logger logging.Logger) *Seeder {
	return &Seeder{db: db}
}

func (s *Seeder) RunSeeders() {
	seeder := NewSeeder(s.db, s.logger)
	seeder.SeedPermissions()
	seeder.SeedRoles()
}

func (s *Seeder) SeedPermissions() {
	permissions := []models.Permission{
		{Action: "view_survey"},
		{Action: "view_survey_results"},
		{Action: "restrict_user"},
		{Action: "view_specific_participant_votes"},
		{Action: "vote"},
		{Action: "edit_survey"},
		{Action: "view_survey_roles"},
		{Action: "assign_and_remove_survey_roles"},
		{Action: "view_survey_reports"},
	}

	for _, perm := range permissions {
		var existingPermission models.Permission
		if err := s.db.GetDb().Where("action = ?", perm.Action).First(&existingPermission).Error; err != nil {
			if err.Error() == "record not found" {
				err := s.db.GetDb().Create(&perm).Error
				if err != nil {
					s.logger.Error(logging.Database, logging.Insert, "create permission error in seeder", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				}
			} else {
				log.Printf("Error checking permission: %v", err)
			}
		}
	}
}

func (s *Seeder) SeedRoles() {
	superAdminRole := models.Role{Name: "SuperAdmin"}
	adminRole := models.Role{Name: "Admin"}
	userRole := models.Role{Name: "User"}

	if err := s.db.GetDb().Where("name = ?", "SuperAdmin").First(&superAdminRole).Error; err != nil {
		if err.Error() == "record not found" {
			err := s.db.GetDb().Create(&superAdminRole).Error
			if err != nil {
				s.logger.Error(logging.Database, logging.Insert, "create super admin error in seeder", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			}
		} else {
			log.Printf("Error checking super admin role: %v", err)
		}
	}

	if err := s.db.GetDb().Where("name = ?", "Admin").First(&adminRole).Error; err != nil {
		if err.Error() == "record not found" {
			err := s.db.GetDb().Create(&adminRole).Error
			if err != nil {
				s.logger.Error(logging.Database, logging.Insert, "create admin error in seeder", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			}
		} else {
			log.Printf("Error checking admin role: %v", err)
		}
	}

	if err := s.db.GetDb().Where("name = ?", "User").First(&userRole).Error; err != nil {
		if err.Error() == "record not found" {
			err := s.db.GetDb().Create(&userRole).Error
			if err != nil {
				s.logger.Error(logging.Database, logging.Insert, "create user error in seeder", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			}
		} else {
			log.Printf("Error checking user role: %v", err)
		}
	}

	var superAdminPermissions []models.Permission
	err := s.db.GetDb().
		Where("action IN ?", []string{"view_survey", "view_survey_results", "restrict_user"}).
		Find(&superAdminPermissions).Error
	if err != nil {
		s.logger.Error(logging.Database, logging.Select, "select permissions in seeder", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}

	for _, perm := range superAdminPermissions {
		if err := s.db.GetDb().Model(&superAdminRole).Association("Permissions").Append(&perm); err != nil {
			s.logger.Error(logging.Database, logging.Insert, "assign permission to super admin error", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		}
	}
}
