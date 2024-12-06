package service

import (
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type IQuestionService interface {
	GetQuestion(c context.Context, id uint) (response *dto.Question, err error)
	GetQuestions(c context.Context, req dto.GetQuestionsRequest) (response []*dto.Question, err error)
	UpdateQuestion(c context.Context, id uint, req dto.QuestionUpdateRequest) (question *dto.Question, err error)
	DeleteQuestion(c context.Context, id uint) error
}

type QuestionService struct {
	conf   *config.Config
	repo   repository.ISurveyRepository
	logger logging.Logger
}

func NewQuestionService(conf *config.Config, repo repository.ISurveyRepository, logger logging.Logger) *QuestionService {
	return &QuestionService{conf: conf, repo: repo, logger: logger}
}

func (q *QuestionService) GetQuestions(c context.Context, req dto.GetQuestionsRequest) (response []*dto.Question, err error) {

	Filters := []*dto.RepositoryFilter{}
	Filter := dto.RepositoryFilter{Field: "survey_id", Operator: "=", Value: strconv.Itoa(int(req.SurveyId))}
	Filters = append(Filters, &Filter)
	questions, err := q.repo.GetQuestions(c, &dto.RepositoryRequest{Filters: Filters, With: "Choices"})
	if err != nil {
		return []*dto.Question{}, err
	}

	return response, util.ConvertTypes(q.logger, questions, &response)
}
func (q *QuestionService) UpdateQuestion(c context.Context, id uint, req dto.QuestionUpdateRequest) (question *dto.Question, err error) {
	mq, err := q.repo.GetQuestionByID(c, id)
	if err != nil {
		return nil, err
	}
	if mq == nil {
		return nil, nil
	}
	mq.Text = req.Text
	mq.HasMultipleChoice = req.HasMultipleChoice
	mq.MediaUrl = req.MediaUrl

	_, err = q.repo.UpdateQuestion(c, mq)
	if err != nil {
		return nil, err
	}

	choices := req.Choices
	if len(choices) > 0 && req.HasMultipleChoice {
		err := q.repo.DeleteQuestionChoices(c, id)
		if err != nil {
			return nil, err
		}
		for _, v := range choices {
			ch := models.Choice{}
			ch.IsCorrect = v.IsCorrect
			ch.LinkedQuestionID = v.LinkedQuestionId
			ch.Text = v.Text
			ch.QuestionID = id
			err := q.repo.CreateChoice(c, &ch)
			if err != nil {
				return nil, err
			}
		}

	}
	return question, util.ConvertTypes(q.logger, mq, &question)

}
func (q *QuestionService) DeleteQuestion(c context.Context, id uint) error {

	err := q.repo.DeleteQuestion(c, id)
	if err != nil {
		return err
	}
	return q.repo.DeleteQuestionChoices(c, id)

}
func (q *QuestionService) GetQuestion(c context.Context, id uint) (response *dto.Question, err error) {
	qu, err := q.repo.GetQuestionByID(c, id)

	if err != nil {
		return nil, err
	}

	if qu == nil {
		return nil, nil
	}
	return response, util.ConvertTypes(q.logger, qu, &response)
}
