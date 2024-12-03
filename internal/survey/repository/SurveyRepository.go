package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"gorm.io/gorm"
)

type ISurveyRepository interface {
	CreateSurvey(ctx context.Context, survey *models.Survey) error
	GetSurveyByID(ctx context.Context, surveyId uint) (*models.Survey, error)
	CreateQuestion(ctx context.Context, question *models.Question) error
	CreateChoice(ctx context.Context, choice *models.Choice) error
	UpdateChoice(ctx context.Context, choice *models.Choice) error
	GetChoiceByTextAndQuestion(ctx context.Context, text string, questionID uint) (*models.Choice, error)
	GetUserParticipationList(ctx context.Context, userId uint, surveyId uint) ([]models.UserSurveyParticipation, error)
	CreateUserParticipation(ctx context.Context, participation *models.UserSurveyParticipation) (*models.UserSurveyParticipation, error)
	UpdateUserParticipation(ctx context.Context, participation *models.UserSurveyParticipation) error
	GetUserParticipation(ctx context.Context, participationId uint) (*models.UserSurveyParticipation, error)
}

type SurveyRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewSurveyRepository(db db.DbService, logger logging.Logger) *SurveyRepository {
	return &SurveyRepository{db: db, logger: logger}
}

func (r *SurveyRepository) GetSurveyByID(ctx context.Context, surveyId uint) (*models.Survey, error) {
	var survey models.Survey

	err := r.db.GetDb().WithContext(ctx).First(&survey, surveyId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &survey, err
}
func (r *SurveyRepository) CreateSurvey(ctx context.Context, survey *models.Survey) error {
	err := r.db.GetDb().WithContext(ctx).Create(&survey).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create survey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	err := r.db.GetDb().WithContext(ctx).Create(&question).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create question error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) CreateChoice(ctx context.Context, choice *models.Choice) error {
	err := r.db.GetDb().WithContext(ctx).Create(&choice).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create choice error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) UpdateChoice(ctx context.Context, choice *models.Choice) error {
	err := r.db.GetDb().WithContext(ctx).Save(choice).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "Get choice by text and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) GetChoiceByTextAndQuestion(ctx context.Context, text string, questionID uint) (*models.Choice, error) {
	var choice models.Choice
	err := r.db.GetDb().WithContext(ctx).Where("text = ? AND question_id = ?", text, questionID).First(&choice).Error

	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "update choice by text and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return &choice, err
}

func (r *SurveyRepository) GetUserParticipationList(ctx context.Context, userId uint, surveyId uint) ([]models.UserSurveyParticipation, error) {
	var userParticipationList []models.UserSurveyParticipation
	err := r.db.GetDb().WithContext(ctx).Where("user_id = ? AND survey_id = ?", userId, surveyId).Find(&userParticipationList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.UserSurveyParticipation{}, nil
		}
		r.logger.Error(logging.Database, logging.Select, "Get user participation list  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return userParticipationList, nil
}

func (r *SurveyRepository) CreateUserParticipation(ctx context.Context, participation *models.UserSurveyParticipation) (*models.UserSurveyParticipation, error) {
	err := r.db.GetDb().WithContext(ctx).Create(&participation).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create choice UserSurveyParticipation in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return participation, err
}

func (r *SurveyRepository) UpdateUserParticipation(ctx context.Context, participation *models.UserSurveyParticipation) error {
	err := r.db.GetDb().WithContext(ctx).Save(participation).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "update UserSurveyParticipation  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) GetUserParticipation(ctx context.Context, participationId uint) (*models.UserSurveyParticipation, error) {
	var p models.UserSurveyParticipation

	err := r.db.GetDb().WithContext(ctx).First(&p, participationId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &p, err
}
