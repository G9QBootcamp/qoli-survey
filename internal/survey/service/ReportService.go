package service

import (
	"fmt"
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IReportService interface {
	GetTotalParticipationPercentage(ctx context.Context, surveyId uint) (uint, error)
	GetCorrectAnswerPercentage(ctx context.Context, surveyId uint) ([]dto.CorrectAnswerPercentageToShow, error)
	SuddenlyFinishedParticipationPercentage(ctx context.Context, surveyId uint) (float64, error)
	GetChoicesByPercentage(ctx context.Context, surveyId uint) ([]dto.QuestionReport, error)
	GetMultipleParticipationCount(ctx context.Context, surveyId uint) ([]dto.ParticipationReport, error)
	GetAverageResponseTime(ctx context.Context, surveyId uint) (float64, error)
	GetResponseDispersionByHour(ctx context.Context, surveyId uint) ([]dto.HourDispersionDTO, error)
}
type ReportService struct {
	conf   *config.Config
	repo   repository.IReportRepository
	logger logging.Logger
}

func NewReportService(conf *config.Config, repo repository.IReportRepository, logger logging.Logger) *ReportService {
	return &ReportService{conf: conf, repo: repo, logger: logger}
}

func (s *ReportService) GetTotalParticipationPercentage(ctx context.Context, surveyId uint) (uint, error) {
	participated, err := s.repo.GetTotalParticipatesForSurvey(ctx, surveyId)
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
	if err != nil {
		return nil, err
	}
	res := make([]dto.CorrectAnswerPercentageToShow, 0)
	for _, q := range qs {
		if q.HasMultipleChoice {
			correctAns, err := s.repo.GetCorrectChoiceByQuestionID(ctx, q.ID)

			if err != nil {
				return nil, err
			}
			if correctAns == nil {
				continue
			}

			totalVotesCount, err := s.repo.GetTotalVotesToQuestionCount(ctx, q.ID)
			if err != nil {
				return nil, err
			}
			if totalVotesCount == 0 {
				res = append(res, dto.CorrectAnswerPercentageToShow{
					QuestionID: q.ID,
					Percentage: "0%",
				})
				continue
			}
			correctAnsCount, err := s.repo.GetGivenAnswerCountByQuestionID(ctx, q.ID, correctAns.Text)
			if err != nil {
				return nil, err
			}
			res = append(res, dto.CorrectAnswerPercentageToShow{
				QuestionID: q.ID,
				Percentage: strconv.Itoa(int(100*(float64(correctAnsCount)/float64(totalVotesCount)))) + "%",
			})
		}
	}
	return res, nil
}

func (s *ReportService) GetMultipleParticipationCount(ctx context.Context, surveyId uint) ([]dto.ParticipationReport, error) {
	users, err := s.repo.GetTotalParticipants(ctx, surveyId)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ParticipationReport, 0)

	for _, user := range users {
		participationCount, err := s.repo.GetParticipationCount(ctx, surveyId, user.ID)
		if err != nil {
			return nil, err
		}

		if participationCount > 1 {
			res = append(res, dto.ParticipationReport{
				UserID: user.ID,
				Count:  participationCount,
			})
		}
	}
	return res, err
}

func (s *ReportService) SuddenlyFinishedParticipationPercentage(ctx context.Context, surveyId uint) (float64, error) {
	totalParticipation, err := s.repo.GetTotalParticipatesForSurvey(ctx, surveyId)
	if err != nil {
		return 0, err
	}
	suddenlyFinished, err := s.repo.GetSuddenlyFinishedParticipatesForSurvey(ctx, surveyId)
	if err != nil {
		return 0, err
	}
	if totalParticipation == 0 {
		return 0, nil
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
				if totalVotesCount == 0 {
					questionReport.ChoiceReport = append(questionReport.ChoiceReport, dto.ChoiceReport{
						ID:         choice.ID,
						Text:       choice.Text,
						Percentage: "0 %",
					})
					continue
				}

				chosenCount, err := s.repo.GetGivenAnswerCountByQuestionID(ctx, q.ID, choice.Text)
				if err != nil {
					return nil, err
				}
				questionReport.ChoiceReport = append(questionReport.ChoiceReport, dto.ChoiceReport{
					ID:         choice.ID,
					Text:       choice.Text,
					Percentage: strconv.Itoa(int(100*float64(chosenCount)/float64(totalVotesCount))) + "%",
				})
			}
			res = append(res, questionReport)
		}
	}
	return res, nil
}

func (s *ReportService) GetAverageResponseTime(ctx context.Context, surveyId uint) (float64, error) {
	return s.repo.GetAverageResponseTime(ctx, surveyId)
}

func (s *ReportService) GetResponseDispersionByHour(ctx context.Context, surveyId uint) ([]dto.HourDispersionDTO, error) {
	dispersionData, err := s.repo.GetResponseDispersionByHour(ctx, surveyId)
	if err != nil {
		return nil, err
	}

	var result []dto.HourDispersionDTO
	for hour, count := range dispersionData {
		result = append(result, dto.HourDispersionDTO{
			Hour:  hour,
			Count: count,
		})
	}

	return result, nil
}
