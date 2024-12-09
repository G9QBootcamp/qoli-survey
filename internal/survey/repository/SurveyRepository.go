package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
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
	GetLastUserParticipation(ctx context.Context, userId uint, surveyId uint) (*models.UserSurveyParticipation, error)
	CreateVote(ctx context.Context, v *models.Vote) (*models.Vote, error)
	UpdateVote(ctx context.Context, v *models.Vote) (*models.Vote, error)
	GetUserSurveyVote(ctx context.Context, user_id uint, question_id uint) (*models.Vote, error)
	UpdateQuestion(c context.Context, m *models.Question) (*models.Question, error)
	DeleteQuestion(c context.Context, id uint) error
	DeleteSurvey(c context.Context, id uint) error
	GetQuestionByID(ctx context.Context, id uint) (*models.Question, error)
	DeleteQuestionChoices(ctx context.Context, questionId uint) error
	GetQuestions(ctx context.Context, req *dto.RepositoryRequest) ([]*models.Question, error)
	GetSurveys(ctx context.Context, req *dto.RepositoryRequest) (questions []*models.Survey, err error)
	DeleteVote(c context.Context, id uint) error
	GetVoteByID(ctx context.Context, id uint) (*models.Vote, error)
	UpdateSurvey(ctx context.Context, survey *models.Survey) error
	GetSurveyVotes(ctx context.Context, id uint) ([]*models.Vote, error)

	CreateOption(ctx context.Context, option *models.SurveyOption) (*models.SurveyOption, error)
	UpdateOption(ctx context.Context, option *models.SurveyOption) error
	DeleteOption(c context.Context, id uint) error
	GetOptionByID(ctx context.Context, id uint) (*models.SurveyOption, error)
	GetOptions(ctx context.Context, req *dto.RepositoryRequest) (options []*models.SurveyOption, err error)
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

	err := r.db.GetDb().WithContext(ctx).Preload("Options").First(&survey, surveyId).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &survey, err
}

func (r *SurveyRepository) GetVoteByID(ctx context.Context, id uint) (*models.Vote, error) {
	var vote models.Vote

	err := r.db.GetDb().WithContext(ctx).First(&vote, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &vote, err
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

func (r *SurveyRepository) GetQuestions(ctx context.Context, req *dto.RepositoryRequest) (questions []*models.Question, err error) {
	return GetRecords[*models.Question](r.db.GetDb(), req)
}
func (r *SurveyRepository) GetSurveys(ctx context.Context, req *dto.RepositoryRequest) (questions []*models.Survey, err error) {
	return GetRecords[*models.Survey](r.db.GetDb(), req)
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

func (r *SurveyRepository) UpdateSurvey(ctx context.Context, survey *models.Survey) error {
	err := r.db.GetDb().WithContext(ctx).Save(survey).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "Get survey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
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

func (r *SurveyRepository) CreateVote(ctx context.Context, v *models.Vote) (*models.Vote, error) {
	err := r.db.GetDb().WithContext(ctx).Create(&v).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create Vote error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return v, err
}
func (r *SurveyRepository) UpdateVote(ctx context.Context, v *models.Vote) (*models.Vote, error) {
	err := r.db.GetDb().WithContext(ctx).Save(v).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "update vote  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return v, err
}
func (r *SurveyRepository) GetUserSurveyVote(ctx context.Context, user_id uint, question_id uint) (*models.Vote, error) {
	var vote models.Vote
	err := r.db.GetDb().WithContext(ctx).Where("voter_id = ? AND question_id = ?", user_id, question_id).First(&vote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error(logging.Database, logging.Select, "get vote by user id and question id error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return &vote, err
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
func (r *SurveyRepository) GetLastUserParticipation(ctx context.Context, userId uint, surveyId uint) (*models.UserSurveyParticipation, error) {
	var userParticipation *models.UserSurveyParticipation
	err := r.db.GetDb().WithContext(ctx).Where("user_id = ? AND survey_id = ?", userId, surveyId).Last(&userParticipation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error(logging.Database, logging.Select, "Get last user participation  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return userParticipation, nil
}

func (q *SurveyRepository) UpdateQuestion(c context.Context, m *models.Question) (*models.Question, error) {
	err := q.db.GetDb().WithContext(c).Save(m).Error
	if err != nil {
		q.logger.Error(logging.Database, logging.Update, "update question  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return m, err
}
func (q *SurveyRepository) DeleteQuestion(c context.Context, id uint) error {
	return q.db.GetDb().WithContext(c).Where("ID = ?", id).Delete(&models.Question{}).Error

}

func (q *SurveyRepository) DeleteSurvey(c context.Context, id uint) error {
	return q.db.GetDb().WithContext(c).Where("ID = ?", id).Delete(&models.Survey{}).Error

}
func (q *SurveyRepository) DeleteVote(c context.Context, id uint) error {
	return q.db.GetDb().WithContext(c).Where("ID = ?", id).Delete(&models.Vote{}).Error

}

func (r *SurveyRepository) GetQuestionByID(ctx context.Context, id uint) (*models.Question, error) {
	var question models.Question

	err := r.db.GetDb().WithContext(ctx).Preload("Choices").First(&question, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &question, err
}

func (r *SurveyRepository) DeleteQuestionChoices(c context.Context, questionId uint) error {
	return r.db.GetDb().WithContext(c).Where("question_id = ?", questionId).Delete(&models.Choice{}).Error

}

func (r *SurveyRepository) CreateOption(ctx context.Context, option *models.SurveyOption) (*models.SurveyOption, error) {
	err := r.db.GetDb().WithContext(ctx).Create(&option).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Insert, "create choice UserSurveyParticipation in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return option, err
}
func (r *SurveyRepository) UpdateOption(ctx context.Context, option *models.SurveyOption) error {
	err := r.db.GetDb().WithContext(ctx).Save(option).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Update, "update vote  error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return err
}
func (r *SurveyRepository) DeleteOption(c context.Context, id uint) error {
	return r.db.GetDb().WithContext(c).Where("ID = ?", id).Delete(&models.SurveyOption{}).Error

}
func (r *SurveyRepository) GetOptionByID(ctx context.Context, id uint) (option *models.SurveyOption, err error) {

	err = r.db.GetDb().WithContext(ctx).First(&option, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return option, err

}
func (r *SurveyRepository) GetOptions(ctx context.Context, req *dto.RepositoryRequest) (options []*models.SurveyOption, err error) {
	return GetRecords[*models.SurveyOption](r.db.GetDb(), req)

}

func (r *SurveyRepository) GetSurveyVotes(ctx context.Context, id uint) ([]*models.Vote, error) {
	var votes []*models.Vote
	err := r.db.GetDb().Preload("Voter").
		Preload("Question.Survey").
		Where("question_id IN (SELECT id FROM questions WHERE survey_id = ?)", id).
		Find(&votes).Error
	return votes, err
}
