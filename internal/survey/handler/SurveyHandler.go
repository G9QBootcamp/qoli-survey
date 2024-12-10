package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
	"github.com/gorilla/websocket"

	notification "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/service"

	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type SurveyHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.ISurveyService
	logger  logging.Logger
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewSurveyHandler(conf *config.Config, db db.DbService, logger logging.Logger, notificationService notification.INotificationService) *SurveyHandler {
	return &SurveyHandler{conf: conf, db: db, service: service.NewSurveyService(conf, repository.NewSurveyRepository(db, logger), logger, notificationService), logger: logger}
}

func (h *SurveyHandler) CreateSurvey(c echo.Context) error {
	var req dto.SurveyCreateRequest

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in create survey api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	for _, question := range req.Questions {
		seen := make(map[string]bool)
		for _, choice := range question.Choices {
			if seen[strings.ToLower(choice.Text)] {
				h.logger.Info(logging.Internal, logging.Api, "validation error in create survey api", map[logging.ExtraKey]interface{}{logging.Service: "SurveyService"})

				return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed: the request has same choices for a question"})
			}
			seen[strings.ToLower(choice.Text)] = true
		}
	}

	req.OwnerID = userID
	survey, err := h.service.CreateSurvey(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, survey)
}

func (h *SurveyHandler) GetSurveys(c echo.Context) error {
	var req dto.SurveysGetRequest

	userID, ok := c.Get("userID").(uint)
	role, _ := c.Get("role").(string)

	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	if role != "SuperAdmin" {
		req.UserId = int(userID)

	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in get surveys api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	response, err := h.service.GetSurveys(c.Request().Context(), req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)

}

func (h *SurveyHandler) GetSurvey(c echo.Context) error {

	survey_id := c.Param("survey_id")
	userID, ok := c.Get("userID").(uint)

	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	survey, err := h.service.GetSurvey(c.Request().Context(), uint(iSurveyId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if survey == nil {
		return c.JSON(http.StatusNotFound, nil)

	}
	return c.JSON(http.StatusOK, survey)

}
func (h *SurveyHandler) DeleteSurvey(c echo.Context) error {
	survey_id := c.Param("survey_id")
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	err = h.service.DeleteSurvey(c.Request().Context(), uint(iSurveyId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}
func (h *SurveyHandler) StartSurvey(c echo.Context) error {

	survey_id := c.Param("survey_id")
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in start survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	survey, err := h.service.GetSurvey(c.Request().Context(), uint(iSurveyId))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Failed to get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if survey == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "survey not found"})

	}

	can, canError := h.service.CanUserParticipateToSurvey(c.Request().Context(), userID, uint(iSurveyId))
	if canError != nil || !can {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": canError.Error()})
	}

	questionsAnswerMap, err := h.service.GetSurveyQuestionsInOrder(c.Request().Context(), survey.SurveyID)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if len(questionsAnswerMap) <= 0 {
		return c.JSON(http.StatusNoContent, map[string]string{"error": "there are no question for this survey"})
	}

	timeLimit := survey.AnswerTimeLimit

	participation, err := h.service.Participate(c.Request().Context(), userID, survey.SurveyID)
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "error in create user participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Failed to upgrade connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}
	defer conn.Close()

	// Channel to signal disconnection
	disconnectSignal := make(chan struct{})

	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(timeLimit)*time.Second)
	defer cancel()

	go h.startTimer(ctx, conn, participation.ID, disconnectSignal)
	h.readAnswers(ctx, conn, participation.ID, userID, questionsAnswerMap, survey.AllowReturn, disconnectSignal)

	return nil
}

func (h *SurveyHandler) startTimer(contextWithTimeout context.Context, conn *websocket.Conn, participationId uint, disconnectSignal chan struct{}) error {

	defer h.service.EndParticipation(contextWithTimeout, participationId)

	select {
	case <-disconnectSignal:
		err := h.service.EndParticipation(context.Background(), participationId)
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "error in ending user survey participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			return err
		}
		return nil
	case <-contextWithTimeout.Done():

		err := h.service.EndParticipation(context.Background(), participationId)
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "error in ending user survey participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			return err
		}
		err = conn.WriteJSON(dto.VoteResponse{Question: nil, Message: "time limit reached"})
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		}
		return conn.Close()
	}

}

func (h *SurveyHandler) readAnswers(c context.Context, conn *websocket.Conn, participationId uint, userId uint, questionsAnswerMap dto.QuestionsAnswerMap, allowReturn bool, disconnectSignal chan struct{}) {
	defer close(disconnectSignal)

	sentQuestions := []*dto.Question{}
	i := 0
	items := questionsAnswerMap[i]
	q := items[dto.NoAnswer]
	err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "answer question:"})
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return
	}
	sentQuestions = append(sentQuestions, q)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error(logging.General, logging.Api, "Unexpected close websocket connection error", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			} else {
				h.logger.Info(logging.General, logging.Api, "Websocket Connection closed", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			}
			return
		}
		req := dto.VoteRequest{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "Unexpected close websocket connection error", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		}

		if req.Operation != dto.CommitOperation && req.Operation != dto.BackOperation {
			err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "invalid operation"})
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return
			}
		}

		if req.Operation == dto.BackOperation && allowReturn {
			if i == 0 {
				err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "this is first question"})
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
					return
				}
				continue
			}

			i--
			sentQuestions = sentQuestions[:len(sentQuestions)-1]
			q = sentQuestions[len(sentQuestions)-1]

			err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "answer question:"})
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return
			}
			continue
		}

		if req.Operation == dto.BackOperation && !allowReturn {
			err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "you are not allowed to return in this survey"})
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return
			}
			continue
		}

		if req.QuestionId != q.ID {
			err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "invalid question id"})
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

			}

			continue
		}

		validChoice := false
		isCorrectAnswer := false
		if len(q.Choices) > 0 && q.HasMultipleChoice {
			for _, v := range q.Choices {
				if v.Text == req.Answer && v.IsCorrect {
					isCorrectAnswer = true
				}
				if v.Text == req.Answer {
					validChoice = true
				}
			}

			if !validChoice {
				err := conn.WriteJSON(dto.VoteResponse{Question: q, Message: "enter valid choice"})
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})

				}
				continue
			}

		}

		h.service.CommitVote(c, models.Vote{VoterID: userId, QuestionID: q.ID, Answer: req.Answer, IsCorrect: isCorrectAnswer})
		i++

		if len(questionsAnswerMap) <= i {
			break
		}

		for {
			if len(questionsAnswerMap) <= i {
				break
			}
			items := questionsAnswerMap[i]
			var e bool
			q, e = items[dto.Answer(req.Answer)]
			if e {
				break
			} else {
				q, e = items[dto.NoAnswer]
				if e {
					break
				} else {
					i++
				}
			}

		}
		if q == nil {
			break
		}
		err = conn.WriteJSON(dto.VoteResponse{Question: q, Message: "answer question:"})
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			return
		}
		sentQuestions = append(sentQuestions, q)

	}
	err = h.service.CommitParticipation(c, participationId)
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "error in committing survey participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return
	}
	err = conn.WriteJSON(dto.VoteResponse{Message: "survey answers committed successfully"})
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return
	}
}

func (h *SurveyHandler) GetUserVotes(c echo.Context) error {
	viewerID, ok := c.Get("userID").(uint)
	if !ok || viewerID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid survey_id"})
	}
	respondentID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user_id"})
	}

	votes, err := h.service.GetVotes(uint(surveyID), viewerID, uint(respondentID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"votes": votes})
}

func (h *SurveyHandler) GetVisibleVoteUsers(c echo.Context) error {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid survey ID"})
	}

	userID := c.Get("user_id").(uint)

	users, err := h.service.GetVisibleVoteUsers(uint(surveyID), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}
