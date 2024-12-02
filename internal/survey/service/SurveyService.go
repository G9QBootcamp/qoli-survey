package service

import (
	"fmt"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type ISurveyService interface {
	CreateSurvey(c context.Context, req dto.SurveyCreateRequest) (*dto.SurveyResponse, error)
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
