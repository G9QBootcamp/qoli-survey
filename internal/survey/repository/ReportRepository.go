package repository

import (
	"context"

	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type IReportRepository interface {
	GetSurveyParticipantsCount(ctx context.Context, surveyId uint) (int64, error)
	GetSurveyParticipantsCountByPermissionId(ctx context.Context, surveyId uint, permissionId uint) (int64, error)
	GetGivenAnswerCount(ctx context.Context, ans string) (int64, error)
	GetTotalVotesToQuestionCount(ctx context.Context, qid uint) (int64, error)
	GetTotalParticipateCountForSurvey(ctx context.Context, surveyId uint) (int64, error)
	GetSuddenlyFinishedParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error)
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
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "get survey participants count error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		count = 0
		return count, err
	}
	return count, nil
}
func (r *ReportRepository) GetSurveyParticipantsCountByPermissionId(ctx context.Context, surveyId uint, permissionId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Raw("select count(usr.user_id)"+
		"where usr.survey_id = ? AND rp.permission_id = ?"+
		"from user_survey_role as usr"+
		"join role as r"+
		"on usr.role_id = r.ID"+
		"join role_permission as rp"+
		"on r.ID = rp.role_id"+
		"group by usr.user_id",
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
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetSurveyParticipantsCountByPermissionId error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return count, err
}

func (r *ReportRepository) GetGivenAnswerCount(ctx context.Context, ans string) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("vote").Where("answer = ?", ans).Count(&count).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetGivenAnswerCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return count, err
}

func (r *ReportRepository) GetTotalVotesToQuestionCount(ctx context.Context, qid uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("vote").Where("question_id = ?", qid).Count(&count).Error

	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetTotalVotesToQuestionCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return count, err
}
func (r *ReportRepository) GetTotalParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participation").Where("survey_id = ?", surveyId).Count(&count).Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetTotalParticipatesForSurvey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return count, err
}
func (r *ReportRepository) GetSuddenlyFinishedParticipatesForSurvey(ctx context.Context, surveyId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participation").
		Where("survey_id = ? AND committed_at = null and ended_at != null", surveyId).
		Count(&count).
		Error
	if err != nil {
		r.logger.Error(logging.Database, logging.Select, "GetSuddenlyFinishedParticipatesForSurvey error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
	}
	return count, err
}
