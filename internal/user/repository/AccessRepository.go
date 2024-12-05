package repository

import (
	"context"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type IAccessRepository interface {
	CreateRole(ctx context.Context, role models.Role) (*models.Role, error)
	CreateUserSurveyRole(ctx context.Context, usr models.UserSurveyRole) (*models.UserSurveyRole, error)
	GetAllPermissions(ctx context.Context) ([]models.Permission, error)
	DeleteUserSurveyRole(ctx context.Context, surveyID uint, userID uint, roleID uint) error
	GetUserRolesForSurvey(ctx context.Context, userID, surveyID uint) ([]models.UserSurveyRole, error)
}

type AccessRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewAccessRepository(db db.DbService, logger logging.Logger) *AccessRepository {
	return &AccessRepository{db: db, logger: logger}
}

func (r *AccessRepository) CreateRole(ctx context.Context, role models.Role) (*models.Role, error) {
	err := r.db.GetDb().WithContext(ctx).Save(&role).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create role error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return &role, err
}

func (r *AccessRepository) CreateUserSurveyRole(ctx context.Context, usr models.UserSurveyRole) (*models.UserSurveyRole, error) {
	err := r.db.GetDb().WithContext(ctx).Create(&usr).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create role error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return &usr, err
}

func (r *AccessRepository) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.GetDb().WithContext(ctx).Find(&permissions).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.FailedToGetPermissions, "create role error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return permissions, err
}

func (r *AccessRepository) DeleteUserSurveyRole(ctx context.Context, surveyID uint, userID uint, roleID uint) error {
	err := r.db.GetDb().WithContext(ctx).
		Where("survey_id = ? AND user_id = ? AND role_id = ?", surveyID, userID, roleID).
		Delete(&models.UserSurveyRole{}).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Delete, "delete role from user survey role error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}

	err = r.db.GetDb().WithContext(ctx).Where("id = ?", roleID).Delete(&models.Role{}).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Delete, "delete role error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}

	return nil
}

func (r *AccessRepository) GetUserRolesForSurvey(ctx context.Context, userID, surveyID uint) ([]models.UserSurveyRole, error) {
	var roles []models.UserSurveyRole
	err := r.db.GetDb().WithContext(ctx).Preload("Role.Permissions").
		Where("user_id = ? AND survey_id = ?", userID, surveyID).Find(&roles).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "get user roles for survey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return roles, err
}
