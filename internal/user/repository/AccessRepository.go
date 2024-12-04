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
	DeleteUserSurveyRole(ctx context.Context, userSurveyRoleId uint) error
	GetUserRolesForSurvey(ctx context.Context, userId, surveyId uint) ([]models.UserSurveyRole, error)
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

func (r *AccessRepository) DeleteUserSurveyRole(ctx context.Context, userSurveyRoleId uint) error {
	return r.db.GetDb().WithContext(ctx).Where("ID = ?", userSurveyRoleId).Delete(&models.UserSurveyRole{}).Error
}

func (r *AccessRepository) GetUserRolesForSurvey(ctx context.Context, userId, surveyId uint) ([]models.UserSurveyRole, error) {
	var roles []models.UserSurveyRole
	err := r.db.GetDb().WithContext(ctx).Where("user_id = ? AND survey_id = ?", userId, surveyId).Find(&roles).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "get user roles for survey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return roles, err
}
