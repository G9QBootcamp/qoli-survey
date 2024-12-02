package repository

import (
	"context"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type ISurveyRepository interface {
	CreateSurvey(ctx context.Context, survey *models.Survey) error
	CreateQuestion(ctx context.Context, question *models.Question) error
	CreateChoice(ctx context.Context, choice *models.Choice) error
	UpdateChoice(ctx context.Context, choice *models.Choice) error
	GetChoiceByTextAndQuestion(ctx context.Context, text string, questionID uint) (*models.Choice, error)
}

type SurveyRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewSurveyRepository(db db.DbService, logger logging.Logger) *SurveyRepository {
	return &SurveyRepository{db: db, logger: logger}
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
		r.logger.Error(logging.Database, logging.Select, "Get choice by text and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}

func (r *SurveyRepository) GetChoiceByTextAndQuestion(ctx context.Context, text string, questionID uint) (*models.Choice, error) {
	var choice models.Choice
	err := r.db.GetDb().WithContext(ctx).Where("text = ? AND question_id = ?", text, questionID).First(&choice).Error

	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "Get choice by text and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return &choice, err
}
