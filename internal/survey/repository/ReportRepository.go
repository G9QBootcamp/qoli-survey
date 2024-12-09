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
	GetAllSurveys(ctx context.Context) ([]models.Survey, error)
	GetAccessibleSurveys(ctx context.Context, userID uint, permission string) ([]models.Survey, error)
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

func (r *ReportRepository) GetParticipationCount(ctx context.Context, surveyId uint, userId uint) (int64, error) {
	var count int64
	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participations").
		Where("survey_id = ? AND user_id = ?", surveyId, userId).
		Count(&count).
		Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetParticipationCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}
	return count, nil
}

func (r *ReportRepository) GetTotalParticipants(ctx context.Context, surveyId uint) ([]userModels.User, error) {
	var users []userModels.User

	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participations").
		Select("DISTINCT users.*").
		Joins("JOIN users ON user_survey_participations.user_id = users.id").
		Where("user_survey_participations.survey_id = ?", surveyId).
		Scan(&users).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetParticipationCount error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return users, nil
}

func (r *ReportRepository) GetAverageResponseTime(ctx context.Context, surveyId uint) (float64, error) {
	var avgResponseTimeInMinutes float64
	var count int64

	err := r.db.GetDb().WithContext(ctx).Table("user_survey_participations").
		Where("survey_id = ?", surveyId).
		Count(&count).Error

	if count > 0 {
		err = r.db.GetDb().WithContext(ctx).Table("user_survey_participations").
			Select("AVG(EXTRACT(EPOCH FROM COALESCE(committed_at, end_at) - start_at) / 60.0) AS avg_response_time_minutes").
			Where("survey_id = ?", surveyId).
			Scan(&avgResponseTimeInMinutes).Error
	} else {
		avgResponseTimeInMinutes = 0
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetAverageResponseTime error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return 0, err
	}

	return avgResponseTimeInMinutes, nil
}

func (r *ReportRepository) GetResponseDispersionByHour(ctx context.Context, surveyId uint) (map[int]int, error) {
	var results []struct {
		Hour  int
		Count int
	}
	dispersionData := make(map[int]int)

	err := r.db.GetDb().Table("user_survey_participations").
		Select("EXTRACT(HOUR FROM start_at) as hour, COUNT(*) as count").
		Where("survey_id = ?", surveyId).
		Group("EXTRACT(HOUR FROM start_at)").
		Scan(&results).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetResponseDispersionByHour error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}

	for _, res := range results {
		dispersionData[res.Hour] = res.Count
	}

	return dispersionData, nil
}

func (r *ReportRepository) GetAllSurveys(ctx context.Context) ([]models.Survey, error) {
	var surveys []models.Survey
	err := r.db.GetDb().WithContext(ctx).Find(&surveys).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetAllSurveys error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return surveys, nil
}

func (r *ReportRepository) GetAccessibleSurveys(ctx context.Context, userID uint, permission string) ([]models.Survey, error) {
	var surveys []models.Survey
	err := r.db.GetDb().Joins("JOIN user_survey_roles usr ON usr.survey_id = surveys.id").
		Joins("JOIN roles r ON r.id = usr.role_id").
		Joins("JOIN role_permissions rp ON rp.role_id = r.id").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("usr.user_id = ? AND p.action = ?", userID, permission).
		Find(&surveys).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error(logging.Database, logging.Select, "GetAccessibleSurveys error in repository ", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}
	return surveys, nil
}
