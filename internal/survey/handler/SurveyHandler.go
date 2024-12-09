package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
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

	req.OwnerID = userID
	survey, err := h.service.CreateSurvey(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, survey)
}

func (h *SurveyHandler) CreateSurveyOption(c echo.Context) error {
	var req dto.SurveyOptionCreateRequest

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	survey_id := c.Param("survey_id")
	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in create option survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in create survey api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	survey, err := h.service.CreateOption(c.Request().Context(), userID, uint(iSurveyId), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, survey)
}

func (h *SurveyHandler) GetSurveyOptions(c echo.Context) error {

	userID, ok := c.Get("userID").(uint)

	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}
	survey_id := c.Param("survey_id")
	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in create option survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	response, err := h.service.GetOptions(c.Request().Context(), dto.SurveyOptionsGetRequest{SurveyId: uint(iSurveyId)})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)

}

func (h *SurveyHandler) UpdateSurveyOption(c echo.Context) error {
	var req dto.SurveyOptionCreateRequest

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	optionId := c.Param("option_id")
	iOptionId, err := strconv.Atoi(optionId)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in update option survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in create survey api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	survey, err := h.service.UpdateOption(c.Request().Context(), uint(iOptionId), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, survey)
}

func (h *SurveyHandler) DeleteSurveyOption(c echo.Context) error {

	userID, ok := c.Get("userID").(uint)

	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}
	optionId := c.Param("option_id")
	iOptionId, err := strconv.Atoi(optionId)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in update option survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	err = h.service.DeleteOption(c.Request().Context(), uint(iOptionId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)

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
func (h *SurveyHandler) DeleteVote(c echo.Context) error {
	survey_id := c.Param("survey_id")

	vote_id := c.Param("vote_id")
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}
	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	survey, err := h.service.GetSurvey(c.Request().Context(), uint(iSurveyId))
	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}
	if survey == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "survey not found"})

	}

	ivote_id, err := strconv.Atoi(vote_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete vote", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid vote id"})
	}

	vote, err := h.service.GetVote(c.Request().Context(), uint(ivote_id))

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete vote", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid vote id"})
	}

	if vote == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "vote not found"})

	}

	role, _ := c.Get("role").(string)

	if vote.VoterID != userID && role != "SuperAdmin" {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "you are not allowed"})

	}

	for _, v := range survey.Options {

		if v.Name == "vote_deletion_limit_hours" {
			hours, err := strconv.Atoi(v.Value)
			if err != nil {
				continue
			}

			if vote.CreatedAt.Add(time.Hour * time.Duration(hours)).Before(time.Now()) {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "you are not allowed because of vote_deletion_limit_hours"})
			}
		}

	}
	err = h.service.DeleteVote(c.Request().Context(), uint(ivote_id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}
func (h *SurveyHandler) UpdateSurvey(c echo.Context) error {
	survey_id := c.Param("survey_id")
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in update survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}

	req := dto.SurveyUpdateRequest{}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in update survey api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	survey, err := h.service.UpdateSurvey(c.Request().Context(), uint(iSurveyId), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, survey)
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

func (h *SurveyHandler) SurveyVotes(c echo.Context) error {

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

	role, _ := c.Get("role").(string)

	for _, v := range survey.Options {

		if v.Name == "votes_visibility" && v.Value == "invisible" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "results of this survey is invisible"})
		}

		if v.Name == "votes_visibility" && v.Value == "admin" && role != "SuperAdmin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "results of this survey is invisible"})
		}

	}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Failed to upgrade connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}
	defer conn.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.logger.Error(logging.General, logging.Api, "Unexpected close websocket connection error", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				} else {
					h.logger.Info(logging.General, logging.Api, "Websocket Connection closed", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				}
				return
			}
		}
	}()
	for {
		select {
		case <-done:
			return nil
		default:
			votes, err := h.service.GetSurveyVotes(c.Request().Context(), uint(iSurveyId))
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "Failed to get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return err
			}
			conn.WriteJSON(votes)
			time.Sleep(2 * time.Second)
		}
	}
}
