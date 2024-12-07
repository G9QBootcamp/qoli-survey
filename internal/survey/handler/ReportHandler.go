package handler

import (
	"net/http"
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IReportService
	logger  logging.Logger
}

func NewReportHandler(conf *config.Config, db db.DbService, logger logging.Logger) *ReportHandler {
	return &ReportHandler{conf: conf, db: db, service: service.NewReportService(conf, repository.NewReportRepository(db, logger), logger), logger: logger}
}

func (h *ReportHandler) GetSurveyReport(c echo.Context) error {
	surveyID, err := strconv.ParseUint(c.Param("survey_id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	var reportResponse dto.ReportResponse

	participationPercentage, err := h.service.GetParticipationPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting participation percentage", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get survey participation percentage"})
	}
	reportResponse.SurveyParticipation = strconv.Itoa(int(participationPercentage)) + "%"

	correctAnswerPercentage, err := h.service.GetCorrectAnswerPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting correct answer percentage", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get correct answer percentage"})
	}
	reportResponse.CorrectAnswers = correctAnswerPercentage

	suddenlyFinishedPercentage, err := h.service.SuddenlyFinishedParticipationPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting suddenly finished participation percentage", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get suddenly finished participation percentage"})
	}
	reportResponse.SuddenlyFinishedParticipation = strconv.FormatFloat(suddenlyFinishedPercentage, 'f', 2, 64) + "%"

	choicesPercentage, err := h.service.GetChoicesByPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting choices percentages", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get choices percentages"})
	}
	reportResponse.ChoicesPercentage = choicesPercentage

	return c.JSON(http.StatusOK, reportResponse)
}
