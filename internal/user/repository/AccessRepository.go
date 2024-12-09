package repository

import (
	"context"
	"errors"
	"time"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type IAccessRepository interface {
	CreateRole(ctx context.Context, role models.Role) (*models.Role, error)
	CreateUserSurveyRole(ctx context.Context, usr models.UserSurveyRole) (*models.UserSurveyRole, error)
	GetAllPermissions(ctx context.Context) ([]models.Permission, error)
	DeleteUserSurveyRole(ctx context.Context, surveyID uint, userID uint, roleID uint) error
	GetUserRolesForSurvey(ctx context.Context, userID, surveyID uint) ([]models.UserSurveyRole, error)
	GetRoleByID(ctx context.Context, roleID uint) (*models.Role, error)
	CreateVoteVisibility(ctx context.Context, request dto.VoteVisibilityCreateRequest) (models.VoteVisibility, error)
	GetVoteVisibilityById(ctx context.Context, id uint) (models.VoteVisibility, error)
	GetVoteVisibilityBySurveyId(ctx context.Context, surveyId uint) ([]models.VoteVisibility, error)
	DeleteVoteVisibilityById(ctx context.Context, id uint) error
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
		Where("user_id = ? AND survey_id = ? AND (expires_at > ? OR expires_at IS NULL)",
			userID, surveyID, time.Now()).Find(&roles).Error

	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "get user roles for survey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return roles, err
}

func (r *AccessRepository) GetRoleByID(ctx context.Context, roleID uint) (*models.Role, error) {
	var role models.Role

	err := r.db.GetDb().WithContext(ctx).First(&role, roleID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
func (r *AccessRepository) CreateVoteVisibility(ctx context.Context, request dto.VoteVisibilityCreateRequest) (models.VoteVisibility, error) {
	vv := models.VoteVisibility{
		SurveyID:     request.SurveyID,
		ViewerID:     request.ViewerID,
		RespondentID: request.RespondentID,
	}

	err := r.db.GetDb().WithContext(ctx).Create(&vv).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "CreateVoteVisibility error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return vv, err
}

func (r *AccessRepository) GetVoteVisibilityById(ctx context.Context, id uint) (models.VoteVisibility, error) {
	var vv models.VoteVisibility
	err := r.db.GetDb().WithContext(ctx).Where("id = ?", id).First(&vv).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "GetVoteVisibilityById error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return vv, err
}

func (r *AccessRepository) GetVoteVisibilityBySurveyId(ctx context.Context, surveyId uint) ([]models.VoteVisibility, error) {
	var vvs []models.VoteVisibility
	err := r.db.GetDb().WithContext(ctx).Where("survey_id = ?", surveyId).Find(&vvs).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetVoteVisibilityBySurveyId error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return vvs, err
}

func (r *AccessRepository) DeleteVoteVisibilityById(ctx context.Context, id uint) error {
	err := r.db.GetDb().WithContext(ctx).Where("id = ?", id).Delete(&models.VoteVisibility{}).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Delete, "DeleteVoteVisibilityById error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}
