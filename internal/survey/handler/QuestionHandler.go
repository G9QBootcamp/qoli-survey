package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type QuestionHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IQuestionService
	logger  logging.Logger
}

func NewQuestionHandler(conf *config.Config, db db.DbService, logger logging.Logger) *QuestionHandler {
	return &QuestionHandler{conf: conf, db: db, service: service.NewQuestionService(conf, repository.NewSurveyRepository(db, logger), logger), logger: logger}
}
func (h *QuestionHandler) GetQuestion(c echo.Context) error {

	question_id := c.Param("question_id")
	userID, ok := c.Get("userID").(uint)

	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iquestion_id, err := strconv.Atoi(question_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in get question", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid question id"})
	}

	q, err := h.service.GetQuestion(c.Request().Context(), uint(iquestion_id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if q == nil {
		return c.JSON(http.StatusNotFound, nil)

	}
	return c.JSON(http.StatusOK, q)

}
func (h *QuestionHandler) DeleteQuestion(c echo.Context) error {
	question_id := c.Param("question_id")
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iQuestion_id, err := strconv.Atoi(question_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete question", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid question id"})
	}

	err = h.service.DeleteQuestion(c.Request().Context(), uint(iQuestion_id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *QuestionHandler) GetQuestions(c echo.Context) error {

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	survey_id := c.Param("survey_id")

	iSurveyId, err := strconv.Atoi(survey_id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid survey id"})
	}
	var req dto.GetQuestionsRequest = dto.GetQuestionsRequest{SurveyId: uint(iSurveyId)}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in get questions api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	response, err := h.service.GetQuestions(c.Request().Context(), req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)

}

func (h *QuestionHandler) UpdateQuestion(c echo.Context) error {
	question_id := c.Param("question_id")

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iQuestion_id, err := strconv.Atoi(question_id)

	var req dto.QuestionUpdateRequest

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in delete question", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid question id"})
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in update question api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	seen := make(map[string]bool)
	for _, choice := range req.Choices {
		if seen[strings.ToLower(choice.Text)] {
			h.logger.Info(logging.Internal, logging.Api, "validation error in update question api", map[logging.ExtraKey]interface{}{logging.Service: "SurveyService"})

			return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed: the request has same choices for a question"})
		}
		seen[strings.ToLower(choice.Text)] = true
	}

	question, err := h.service.UpdateQuestion(c.Request().Context(), uint(iQuestion_id), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, question)
}
