package repository

import (
	"context"
	"errors"

	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	userModels "github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"gorm.io/gorm"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type IReportRepository interface {
	GetSurveyParticipantsCount(ctx context.Context, surveyId uint) (int64, error)
	GetTotalParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error)
	GetSurveyParticipantsCountByPermissionId(ctx context.Context, surveyId uint, permissionId uint) (int64, error)
	GetTotalVotesToQuestionCount(ctx context.Context, qid uint) (int64, error)
	GetSuddenlyFinishedParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error)
	GetQuestionsBySurveyID(ctx context.Context, sid uint) ([]models.Question, error)
	GetCorrectChoiceByQuestionID(ctx context.Context, qid uint) (*models.Choice, error)
	GetChoicesByQuestionID(ctx context.Context, qid uint) ([]models.Choice, error)
	GetGivenAnswerCountByQuestionID(ctx context.Context, qid uint, answer string) (int64, error)
	GetParticipationCount(ctx context.Context, surveyId uint, userId uint) (int64, error)
	GetTotalParticipants(ctx context.Context, surveyId uint) ([]userModels.User, error)
	GetAverageResponseTime(ctx context.Context, surveyId uint) (float64, error)
	GetResponseDispersionByHour(ctx context.Context, surveyId uint) (map[int]int, error)
}

type ReportRepository struct {
	db     db.DbService
	logger logging.Logger
}

func NewReportRepository(db db.DbService, logger logging.Logger) *ReportRepository {
	return &ReportRepository{db: db, logger: logger}
}

func (r *ReportRepository) GetSurveyParticipantsCount(ctx context.Context, surveyId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Where("survey_id = ?", surveyId).Group("user_id").Count(&count).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "get survey participants count error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}
func (r *ReportRepository) GetSurveyParticipantsCountByPermissionId(ctx context.Context, surveyId uint, permissionId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Raw("select count(usr.user_id) "+
		"from user_survey_roles as usr "+
		"join roles as r "+
		"on usr.role_id = r.ID "+
		"join role_permissions as rp "+
		"on r.ID = rp.role_id "+
		"where usr.survey_id = ? AND rp.permission_id = ? ",
		surveyId,
		permissionId).Scan(&count).Error
	/*
		select count(usr.user_id)
		where usr.survey_id = ? AND rp.permission_id = ?
		from user_survey_role as usr
		join role as r
		on usr.role_id = r.ID
		join role_permission as rp
		on r.ID = rp.role_id
		group by usr.user_id

	*/
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetSurveyParticipantsCountByPermissionId error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetGivenAnswerCount(ctx context.Context, ans string) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("votes").Where("answer = ?", ans).Count(&count).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetGivenAnswerCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetTotalVotesToQuestionCount(ctx context.Context, qid uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("votes").Where("question_id = ?", qid).Count(&count).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetTotalVotesToQuestionCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetTotalParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participations").Where("survey_id = ?", surveyId).Count(&count).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetTotalParticipatesForSurvey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetSuddenlyFinishedParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participations").
		Where("survey_id = ? AND committed_at IS NULL AND end_at IS NOT NULL", surveyId).
		Count(&count).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetSuddenlyFinishedParticipatesForSurvey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetQuestionsBySurveyID(ctx context.Context, sid uint) ([]models.Question, error) {
	var questions []models.Question
	err := r.db.GetDb().WithContext(ctx).Where("survey_id = ?", sid).Find(&questions).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetQuestionsBySurveyID error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return questions, nil
}

func (r *ReportRepository) GetCorrectChoiceByQuestionID(ctx context.Context, qid uint) (*models.Choice, error) {
	var choice models.Choice
	err := r.db.GetDb().WithContext(ctx).Where("question_id = ? AND is_correct = true", qid).First(&choice).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetCorrectChoiceByQuestionID error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return &choice, nil
}

func (r *ReportRepository) GetChoicesByQuestionID(ctx context.Context, qid uint) ([]models.Choice, error) {
	var choices []models.Choice
	err := r.db.GetDb().WithContext(ctx).Where("question_id = ?", qid).Find(&choices).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetChoicesByQuestionID error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return choices, nil
}

func (r *ReportRepository) GetGivenAnswerCountByQuestionID(ctx context.Context, qid uint, answer string) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("votes").Where("question_id = ? AND answer = ?", qid, answer).Count(&count).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetGivenAnswerCountByQuestionID error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}
