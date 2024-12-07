package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
	"strconv"
)

type IReportService interface {
	GetParticipationPercentage(ctx context.Context, surveyId uint) (uint, error)
	GetCorrectAnswerPercentage(ctx context.Context, surveyId uint) ([]dto.CorrectAnswerPercentageToShow, error)
	SuddenlyFinishedParticipationPercentage(ctx context.Context, surveyId uint) (float64, error)
	ChoicesPercentages(ctx context.Context, surveyId uint) ([]string, error)
	GetChoicesByPercentage(ctx context.Context, surveyId uint) ([]dto.QuestionReport, error)
}
type ReportService struct {
	conf   *config.Config
	repo   repository.IReportRepository
	logger logging.Logger
}

func NewReportService(conf *config.Config, repo repository.IReportRepository, logger logging.Logger) *ReportService {
	return &ReportService{conf: conf, repo: repo, logger: logger}
}

func (s *ReportService) GetParticipationPercentage(ctx context.Context, surveyId uint) (uint, error) {
	participated, err := s.repo.GetTotalParticipateCountForSurvey(ctx, surveyId)
	if err != nil {
		return 0, err
	}
	allowed, err := s.repo.GetSurveyParticipantsCountByPermissionId(ctx, surveyId, 1)
	if err != nil {
		return 0, err
	}
	return uint(100 * float64(participated) / float64(allowed)), nil
}
func (s *ReportService) GetCorrectAnswerPercentage(ctx context.Context, surveyId uint) ([]dto.CorrectAnswerPercentageToShow, error) {
	qs, err := s.repo.GetQuestionsBySurveyID(ctx, surveyId)
	res := make([]dto.CorrectAnswerPercentageToShow, 0)
	if err != nil {
		return nil, err
	}
	for _, q := range qs {
		if q.HasMultipleChoice {
			correctAns, err := s.repo.GetCorrectChoiceByQuestionID(ctx, q.ID)
			if err != nil {
				return nil, err
			}
			totalVotesCount, err := s.repo.GetTotalVotesToQuestionCount(ctx, q.ID)
			if err != nil {
				return nil, err
			}
			correctAnsCount, err := s.repo.GetGivenAnswerCountByQuestionID(ctx, q.ID, strconv.Itoa(int(correctAns.ID)))
			if err != nil {
				return nil, err
			}
			res = append(res, dto.CorrectAnswerPercentageToShow{
				QuestionID:       q.ID,
				HasCorrectAnswer: true,
				Percentage:       100 * (float64(correctAnsCount) / float64(totalVotesCount)),
			})
		} else {
			res = append(res, dto.CorrectAnswerPercentageToShow{
				QuestionID:       q.ID,
				HasCorrectAnswer: false,
			})
		}
	}
	return res, nil
}

func (s *ReportService) SuddenlyFinishedParticipationPercentage(ctx context.Context, surveyId uint) (float64, error) {
	totalParticipation, err := s.repo.GetTotalParticipateCountForSurvey(ctx, surveyId)
	if err != nil {
		return 0, err
	}
	suddenlyFinished, err := s.repo.GetSurveyParticipantsCount(ctx, surveyId)
	if err != nil {
		return 0, err
	}
	return 100 * float64(suddenlyFinished) / float64(totalParticipation), nil
}

func (s *ReportService) GetChoicesByPercentage(ctx context.Context, surveyId uint) ([]dto.QuestionReport, error) {
	qs, err := s.repo.GetQuestionsBySurveyID(ctx, surveyId)
	res := make([]dto.QuestionReport, 0)
	if err != nil {
		return nil, err
	}
	for _, q := range qs {
		if q.HasMultipleChoice {
			questionReport := dto.QuestionReport{
				QuestionID:   q.ID,
				ChoiceReport: make([]dto.ChoiceReport, 0),
			}
			choices, err := s.repo.GetChoicesByQuestionID(ctx, q.ID)
			if err != nil {
				return nil, err
			}
			totalVotesCount, err := s.repo.GetTotalVotesToQuestionCount(ctx, q.ID)
			if err != nil {
				return nil, err
			}
			for _, choice := range choices {
				chosenCount, err := s.repo.GetGivenAnswerCountByQuestionID(ctx, q.ID, strconv.Itoa(int(choice.ID)))
				if err != nil {
					return nil, err
				}
				questionReport.ChoiceReport = append(questionReport.ChoiceReport, dto.ChoiceReport{
					ID:         choice.ID,
					Percentage: 100 * float64(chosenCount) / float64(totalVotesCount),
				})
			}
			res = append(res, questionReport)
		}
	}
	return res, nil
}
