package service

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	notification "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/util"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"golang.org/x/net/context"
)

type ISurveyService interface {
	CreateSurvey(c context.Context, req dto.SurveyCreateRequest) (*dto.SurveyResponse, error)
	UpdateSurvey(c context.Context, id uint, req dto.SurveyUpdateRequest) (*dto.SurveyResponse, error)
	GetSurvey(c context.Context, id uint) (*dto.SurveyResponse, error)
	GetVote(c context.Context, id uint) (*dto.GetVoteResponse, error)
	GetSurveys(c context.Context, req dto.SurveysGetRequest) ([]*dto.SurveyResponse, error)
	DeleteSurvey(c context.Context, id uint) error
	CanUserParticipateToSurvey(c context.Context, userId uint, surveyId uint) (bool, error)
	Participate(c context.Context, userId uint, surveyId uint) (*dto.UserSurveyParticipationResponse, error)
	EndParticipation(c context.Context, participationId uint) error
	CommitParticipation(c context.Context, participationId uint) error
	CanUserVoteOnSurvey(c context.Context, userId uint, surveyId uint) (bool, error)
	CommitVote(c context.Context, vote models.Vote) error
	DeleteVote(c context.Context, id uint) error
	GetSurveyQuestionsInOrder(c context.Context, surveyId uint) (questionsAnswerMap dto.QuestionsAnswerMap, err error)
	GetVotes(surveyID, viewerID, respondentID uint) ([]map[string]interface{}, error)
	GetVisibleVoteUsers(surveyID, viewerID uint) ([]map[string]interface{}, error)
	GetSurveyVotes(c context.Context, surveyId uint) ([]dto.GetVoteResponse, error)

	CreateOption(c context.Context, userId uint, surveyId uint, req dto.SurveyOptionCreateRequest) (*dto.SurveyOptionResponse, error)
	UpdateOption(c context.Context, id uint, req dto.SurveyOptionCreateRequest) (*dto.SurveyOptionResponse, error)
	DeleteOption(c context.Context, id uint) error
	GetOptions(c context.Context, req dto.SurveyOptionsGetRequest) (response []*dto.SurveyOptionResponse, err error)

	UploadMedia(fileHeader *multipart.FileHeader) (string, error)
}
type SurveyService struct {
	conf                *config.Config
	repo                repository.ISurveyRepository
	logger              logging.Logger
	notificationService notification.INotificationService
}

func NewSurveyService(conf *config.Config, repo repository.ISurveyRepository, logger logging.Logger, notificationService notification.INotificationService) *SurveyService {
	return &SurveyService{conf: conf, repo: repo, logger: logger, notificationService: notificationService}
}

func (s *SurveyService) UpdateSurvey(c context.Context, id uint, req dto.SurveyUpdateRequest) (response *dto.SurveyResponse, err error) {
	survey, err := s.repo.GetSurveyByID(c, id)
	if err != nil {
		return nil, err
	}

	survey.Title = req.Title
	survey.StartTime = req.StartTime
	survey.EndTime = req.EndTime
	survey.AllowReturn = req.AllowReturn
	survey.ParticipationLimit = req.ParticipationLimit
	survey.AnswerTimeLimit = req.AnswerTimeLimit
	survey.IsSequential = req.IsSequential

	err = s.repo.UpdateSurvey(c, survey)
	if err != nil {
		return nil, err
	}

	return response, util.ConvertTypes(s.logger, survey, &response)

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

			}
		}

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

func (s *SurveyService) DeleteSurvey(c context.Context, id uint) error {

	survey, err := s.repo.GetSurveyByID(c, id)
	if err != nil {
		return err
	}
	err = s.repo.DeleteSurvey(c, survey.ID)
	if err != nil {
		return err
	}

	_, err = s.notificationService.Notify(c, survey.OwnerID, "your survey with name: "+survey.Title+" removed")
	return err
}
func (s *SurveyService) DeleteVote(c context.Context, id uint) error {

	vote, err := s.GetVote(c, id)

	if err != nil {
		return err
	}
	if vote == nil {
		return errors.New("vote not found")
	}

	err = s.repo.DeleteVote(c, vote.ID)
	if err != nil {
		return err
	}
	_, err = s.notificationService.Notify(c, vote.VoterID, fmt.Sprintf("your vote with this answer removed: %s", vote.Answer))
	return err
}
func (s *SurveyService) GetSurveys(c context.Context, req dto.SurveysGetRequest) (response []*dto.SurveyResponse, err error) {
	limit := 10
	offset := 0
	if req.Page > 0 {
		offset = limit * (req.Page - 1)
	}

	filter := dto.RepositoryFilter{Field: "title", Operator: "LIKE", Value: req.Title}
	filters := []*dto.RepositoryFilter{&filter}
	if req.UserId > 0 {
		filters = append(filters, &dto.RepositoryFilter{Field: "owner_id", Operator: "=", Value: strconv.Itoa(req.UserId)})
	}

	sort := dto.RepositorySort{Field: "created_at", SortType: "desc"}
	repo_req := dto.RepositoryRequest{Limit: uint(limit), Offset: uint(offset), Filters: filters, Sorts: []*dto.RepositorySort{&sort}}

	surveys, err := s.repo.GetSurveys(c, &repo_req)
	if err != nil {
		return []*dto.SurveyResponse{}, err
	}

	return response, util.ConvertTypes(s.logger, surveys, &response)
}

func (s *SurveyService) GetSurvey(c context.Context, id uint) (*dto.SurveyResponse, error) {

	survey, err := s.repo.GetSurveyByID(c, id)

	if err != nil {
		return nil, err
	}

	if survey == nil {
		return nil, nil
	}
	sResponse := dto.SurveyResponse{}

	err = util.ConvertTypes(s.logger, survey, &sResponse)

	if err != nil {
		return nil, err
	}

	return &sResponse, nil

}
func (s *SurveyService) GetVote(c context.Context, id uint) (*dto.GetVoteResponse, error) {
	vote, err := s.repo.GetVoteByID(c, id)

	if err != nil {
		return nil, err
	}

	if vote == nil {
		return nil, nil
	}
	sResponse := dto.GetVoteResponse{}

	err = util.ConvertTypes(s.logger, vote, &sResponse)

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

func (s *SurveyService) CommitParticipation(c context.Context, participationId uint) error {

	pr, err := s.repo.GetUserParticipation(c, participationId)
	if err != nil {
		s.logger.Error(logging.Internal, logging.FailedToGetParticipation, "error in get user participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

		return err
	}
	now := time.Now()
	pr.CommittedAt = &now
	return s.repo.UpdateUserParticipation(c, pr)
}

func (s *SurveyService) CanUserVoteOnSurvey(c context.Context, userId uint, surveyId uint) (bool, error) {
	participation, err := s.repo.GetLastUserParticipation(c, userId, surveyId)
	if err != nil {
		return false, err
	}

	if participation == nil {
		return false, errors.New("you have not started survey yet")
	}

	if !participation.StartAt.IsZero() && participation.CommittedAt == nil && participation.EndAt == nil {
		return false, errors.New("you have not started survey yet")
	}
	return true, nil
}

func (s *SurveyService) CommitVote(c context.Context, vote models.Vote) error {

	v, err := s.repo.GetUserSurveyVote(c, vote.VoterID, vote.QuestionID)
	if err != nil {
		return err
	}
	if v != nil {
		vote.ID = v.ID
		_, err := s.repo.UpdateVote(c, &vote)
		return err
	}
	_, err = s.repo.CreateVote(c, &vote)
	return err

}

func (s *SurveyService) GetSurveyQuestionsInOrder(c context.Context, surveyId uint) (questionsAnswerMap dto.QuestionsAnswerMap, err error) {

	survey, err := s.repo.GetSurveyByID(c, surveyId)

	if err != nil {

		return nil, err
	}

	if survey == nil {
		return dto.QuestionsAnswerMap{}, nil
	}

	filter := dto.RepositoryFilter{Field: "survey_id", Operator: "=", Value: strconv.Itoa(int(surveyId))}
	sort := dto.RepositorySort{Field: "\"order\"", SortType: "asc"}

	questions, err := s.repo.GetQuestions(c,
		&dto.RepositoryRequest{
			Filters: []*dto.RepositoryFilter{&filter},
			Sorts:   []*dto.RepositorySort{&sort},
			With:    "Choices"})

	if err != nil {
		return dto.QuestionsAnswerMap{}, err
	}

	if len(questions) < 1 {
		return dto.QuestionsAnswerMap{}, err
	}

	if !survey.IsSequential {
		questions = util.ShuffleSlice(questions)
	}

	list := dto.QuestionList{}
	err = util.ConvertTypes(s.logger, questions, &list)
	if err != nil {
		return dto.QuestionsAnswerMap{}, err
	}

	mapQuestions := list.ToMap()
	savedQuestionIds := map[uint]bool{}

	tempQuestionAnswer := map[dto.Answer]*dto.Question{}
	for _, v := range list {

		_, e := savedQuestionIds[v.ID]
		if e {
			continue
		}
		questionsAnswerMap = append(questionsAnswerMap, map[dto.Answer]*dto.Question{dto.NoAnswer: v})
		savedQuestionIds[v.ID] = true
		if len(v.Choices) > 0 {
			for _, z := range v.Choices {
				if z.LinkedQuestionID > 0 {
					qid, e := mapQuestions[z.LinkedQuestionID]
					if e {
						tempQuestionAnswer[dto.Answer(z.Text)] = qid
						savedQuestionIds[qid.ID] = true

					}
				}

			}
		}

		if len(tempQuestionAnswer) > 0 {
			questionsAnswerMap = append(questionsAnswerMap, tempQuestionAnswer)
			tempQuestionAnswer = map[dto.Answer]*dto.Question{}
		}
	}

	return questionsAnswerMap, nil
}

func (s *SurveyService) GetVotes(surveyID, viewerID, respondentID uint) ([]map[string]interface{}, error) {
	hasPermission, err := s.repo.CheckVoteVisibility(surveyID, viewerID, respondentID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("viewer does not have permission to view respondent's votes")
	}

	votes, err := s.repo.GetVotes(surveyID, respondentID)
	if err != nil {
		return nil, err
	}

	response := make([]map[string]interface{}, len(votes))
	for i, vote := range votes {
		response[i] = map[string]interface{}{
			"question_id": vote.QuestionID,
			"answer":      vote.Answer,
		}
	}

	return response, nil
}

func (s *SurveyService) GetVisibleVoteUsers(surveyID, viewerID uint) ([]map[string]interface{}, error) {
	users, err := s.repo.GetVisibleVoteUsers(surveyID, viewerID)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *SurveyService) CreateOption(c context.Context, userId uint, surveyId uint, req dto.SurveyOptionCreateRequest) (response *dto.SurveyOptionResponse, err error) {
	option, err := s.repo.CreateOption(c, &models.SurveyOption{SurveyId: surveyId, Name: req.Name, Value: req.Value, UserId: userId})
	if err != nil {
		return nil, err
	}
	return response, util.ConvertTypes(s.logger, option, &response)

}
func (s *SurveyService) UpdateOption(c context.Context, id uint, req dto.SurveyOptionCreateRequest) (response *dto.SurveyOptionResponse, err error) {
	option, err := s.repo.GetOptionByID(c, id)

	if err != nil {
		return nil, err
	}
	if option == nil {
		return nil, nil
	}
	option.Name = req.Name
	option.Value = req.Value

	err = s.repo.UpdateOption(c, option)

	if err != nil {
		return nil, err
	}
	return response, util.ConvertTypes(s.logger, option, &response)

}
func (s *SurveyService) DeleteOption(c context.Context, id uint) error {
	return s.repo.DeleteOption(c, id)
}
func (s *SurveyService) GetOptions(c context.Context, req dto.SurveyOptionsGetRequest) (response []*dto.SurveyOptionResponse, err error) {
	filters := []*dto.RepositoryFilter{}
	if req.SurveyId > 0 {
		filters = append(filters, &dto.RepositoryFilter{Field: "survey_id", Operator: "=", Value: strconv.Itoa(int(req.SurveyId))})
	}
	if req.Name != "" {
		filters = append(filters, &dto.RepositoryFilter{Field: "survey_id", Operator: "=", Value: req.Name})

	}
	options, err := s.repo.GetOptions(c, &dto.RepositoryRequest{Filters: filters})
	if err != nil {
		return []*dto.SurveyOptionResponse{}, err
	}
	return response, util.ConvertTypes(s.logger, options, &response)
}

func (s *SurveyService) GetSurveyVotes(c context.Context, surveyId uint) (response []dto.GetVoteResponse, err error) {
	votes, err := s.repo.GetSurveyVotes(c, surveyId)
	if err != nil {
		return []dto.GetVoteResponse{}, err
	}

	return response, util.ConvertTypes(s.logger, votes, &response)
}

func (s *SurveyService) UploadMedia(fileHeader *multipart.FileHeader) (string, error) {
	allowedExtensions := map[string]bool{".jpg": true, ".png": true, ".mp4": true, ".mp3": true}
	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[fileExt] {
		return "", fmt.Errorf("invalid file type: %s", fileExt)
	}

	// Read file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		return "", err
	}

	return s.repo.SaveFile(fileHeader.Filename, buf.Bytes())
}
