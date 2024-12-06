package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type ISurveyService interface {
	CreateSurvey(c context.Context, req dto.SurveyCreateRequest) (*dto.SurveyResponse, error)
	GetSurvey(c context.Context, id uint) (*dto.SurveyResponse, error)
	CanUserParticipateToSurvey(c context.Context, userId uint, surveyId uint) (bool, error)
	Participate(c context.Context, userId uint, surveyId uint) (*dto.UserSurveyParticipationResponse, error)
	EndParticipation(c context.Context, participationId uint) error
}
type SurveyService struct {
	conf   *config.Config
	repo   repository.ISurveyRepository
	logger logging.Logger
}

func New(conf *config.Config, repo repository.ISurveyRepository, logger logging.Logger) *SurveyService {
	return &SurveyService{conf: conf, repo: repo, logger: logger}
}

func (s *SurveyService) CreateSurvey(c context.Context, req dto.SurveyCreateRequest) (*dto.SurveyResponse, error) {
	survey := models.Survey{
		Title:              req.Title,
		OwnerID:            req.OwnerID,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		IsSequential:       req.IsSequential,
		AllowReturn:        req.AllowReturn,
		ParticipationLimit: req.ParticipationLimit,
		AnswerTimeLimit:    req.AnswerTimeLimit,
	}

	if err := s.repo.CreateSurvey(c, &survey); err != nil {
		return nil, err
	}

	surveyResponseDTO := &dto.SurveyResponse{
		SurveyID:           survey.ID,
		Title:              survey.Title,
		StartTime:          survey.StartTime.Format("2006-01-02 15:04:05"), // Format as string
		EndTime:            survey.EndTime.Format("2006-01-02 15:04:05"),   // Format as string
		IsSequential:       survey.IsSequential,
		AllowReturn:        survey.AllowReturn,
		ParticipationLimit: survey.ParticipationLimit,
		AnswerTimeLimit:    survey.AnswerTimeLimit,
	}

	questionMap := make(map[string]*models.Question)
	questionOrder := 1
	for _, questionReq := range req.Questions {
		question := models.Question{
			SurveyID:          survey.ID,
			Text:              questionReq.Text,
			HasMultipleChoice: questionReq.HasMultipleChoice,
			MediaUrl:          questionReq.MediaUrl,
		}

		if survey.IsSequential {
			question.Order = questionOrder
			questionOrder++
		}

		if err := s.repo.CreateQuestion(c, &question); err != nil {
			return nil, err
		}

		questionDTO := dto.Question{
			ID:                question.ID,
			Text:              question.Text,
			HasMultipleChoice: question.HasMultipleChoice,
			MediaUrl:          question.MediaUrl,
		}

		if question.HasMultipleChoice {
			for _, choiceReq := range questionReq.Choices {
				choice := models.Choice{
					QuestionID: question.ID,
					Text:       choiceReq.Text,
					IsCorrect:  choiceReq.IsCorrect,
				}

				if err := s.repo.CreateChoice(c, &choice); err != nil {
					return nil, err
				}

				choiceDTO := dto.Choice{
					ID:        choice.ID,
					Text:      choice.Text,
					IsCorrect: choice.IsCorrect,
				}

				questionDTO.Choices = append(questionDTO.Choices, choiceDTO)
			}
		}

		surveyResponseDTO.Questions = append(surveyResponseDTO.Questions, questionDTO)

		questionMap[question.Text] = &question
	}

	for _, q := range req.Questions {
		if q.Condition.QuestionText != "" && q.Condition.Answer != "" {
			condition := q.Condition
			targetQuestion := questionMap[q.Text]

			conditionalQuestion, ok := questionMap[condition.QuestionText]
			if !ok {
				return nil, fmt.Errorf("condition question '%s' not found", condition.QuestionText)
			}

			choice, err := s.repo.GetChoiceByTextAndQuestion(c, condition.Answer, conditionalQuestion.ID)
			if err != nil {
				return nil, err
			}

			choice.LinkedQuestionID = targetQuestion.ID
			if err := s.repo.UpdateChoice(c, choice); err != nil {
				return nil, err
			}
		}
	}

	return surveyResponseDTO, nil
}
func (s *SurveyService) GetSurvey(c context.Context, id uint) (*dto.SurveyResponse, error) {

	survey, err := s.repo.GetSurveyByID(c, id)

	if err != nil {
		return nil, err
	}

	sResponse := dto.SurveyResponse{}

	err = util.ConvertTypes(s.logger, survey, &sResponse)

	if err != nil {
		return nil, err
	}

	return &sResponse, nil

}
func (s *SurveyService) CanUserParticipateToSurvey(c context.Context, userId uint, surveyId uint) (bool, error) {
	userParticipationList, err := s.repo.GetUserParticipationList(c, userId, surveyId)
	if err != nil {
		return false, err
	}
	survey, err := s.repo.GetSurveyByID(c, surveyId)
	if err != nil {
		return false, err
	}
	if survey == nil {
		return false, errors.New("survey does not exists")
	}
	if len(userParticipationList) >= survey.ParticipationLimit {
		return false, errors.New("user participation limit reached ")
	}
	if !time.Now().After(survey.StartTime) {
		return false, errors.New("its not time to start the questionnaire")
	}
	if !time.Now().Before(survey.EndTime) {
		return false, errors.New("questionnaire time ended before")
	}

	for _, v := range userParticipationList {
		if !v.StartAt.IsZero() && v.CommittedAt == nil && v.EndAt == nil {
			return false, errors.New("user participation in this survey has not ended")
		}
	}
	return true, nil

}
func (s *SurveyService) Participate(c context.Context, userId uint, surveyId uint) (*dto.UserSurveyParticipationResponse, error) {
	p, err := s.repo.CreateUserParticipation(c, &models.UserSurveyParticipation{UserId: userId, SurveyID: surveyId, StartAt: time.Now()})

	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToCreateParticipation, "error in participation user to survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return nil, err
	}

	response := dto.UserSurveyParticipationResponse{}

	err = util.ConvertTypes(s.logger, p, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *SurveyService) EndParticipation(c context.Context, participationId uint) error {

	pr, err := s.repo.GetUserParticipation(c, participationId)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGetParticipation, "error in get user participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

		return err
	}
	now := time.Now()
	pr.EndAt = &now
	return s.repo.UpdateUserParticipation(c, pr)
}
