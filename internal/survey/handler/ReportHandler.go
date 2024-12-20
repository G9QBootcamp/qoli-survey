package handler

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/models"
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

	participationPercentage, err := h.service.GetTotalParticipationPercentage(c.Request().Context(), uint(surveyID))
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

	GetMultipleParticipationCount, err := h.service.GetMultipleParticipationCount(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting participation count of each user", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get participation count of each user"})
	}
	reportResponse.MultipleParticipationCount = GetMultipleParticipationCount

	suddenlyFinishedPercentage, err := h.service.SuddenlyFinishedParticipationPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting suddenly finished participation percentage", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get suddenly finished participation percentage"})
	}
	reportResponse.SuddenlyFinishedParticipation = strconv.Itoa(int(suddenlyFinishedPercentage)) + "%"

	choicesPercentage, err := h.service.GetChoicesByPercentage(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting choices percentages", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get choices percentages"})
	}
	reportResponse.ChoicesPercentage = choicesPercentage

	averageResponseTime, err := h.service.GetAverageResponseTime(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting average response time", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get average response time"})
	}
	reportResponse.AverageResponseTime = strconv.Itoa(int(averageResponseTime)) + " minutes"

	dispersionResponseByHour, err := h.service.GetResponseDispersionByHour(c.Request().Context(), uint(surveyID))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error getting response dispersion by hour", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get response dispersion by hour"})
	}
	reportResponse.DispersionResponseByHour = dispersionResponseByHour

	return c.JSON(http.StatusOK, reportResponse)
}

func (h *ReportHandler) GenerateAllSurveysReport(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	var surveys []models.Survey
	var err error

	if c.Get("role") == "SuperAdmin" {
		surveys, err = h.service.GetAllSurveys(c.Request().Context())
	} else {
		surveys, err = h.service.GetAccessibleSurveys(c.Request().Context(), userID, "view_survey_reports")
	}

	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Error fetching surveys", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get surveys"})
	}

	var buffer bytes.Buffer
	csvWriter := csv.NewWriter(&buffer)

	headers := []string{
		"Survey ID",
		"Survey Title",
		"Survey Participation",
		"Correct Answers (%)",
		"Multiple Participation Count",
		"Suddenly Finished Participation (%)",
		"Choices Percentage",
		"Average Response Time (minutes)",
		"Dispersion Response by Hour",
	}
	if err := csvWriter.Write(headers); err != nil {
		h.logger.Error(logging.General, logging.Api, "Error writing CSV headers", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write CSV headers"})
	}

	for _, survey := range surveys {
		report, err := h.service.GetSurveyReport(c.Request().Context(), survey.ID) // Fetch full report for each survey
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "Error generating survey report", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), "survey_id": survey.ID})
			continue
		}

		correctAnswers := formatCorrectAnswers(report.CorrectAnswers)
		choicesPercentage := formatChoicesPercentage(report.ChoicesPercentage)
		dispersionByHour := formatDispersionByHour(report.DispersionResponseByHour)

		row := []string{
			strconv.Itoa(int(survey.ID)),
			survey.Title,
			report.SurveyParticipation,
			correctAnswers,
			formatParticipationReport(report.MultipleParticipationCount),
			report.SuddenlyFinishedParticipation,
			choicesPercentage,
			report.AverageResponseTime,
			dispersionByHour,
		}

		if err := csvWriter.Write(row); err != nil {
			h.logger.Error(logging.General, logging.Api, "Error writing CSV row", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), "survey_id": survey.ID})
		}
	}

	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		h.logger.Error(logging.General, logging.Api, "Error finalizing CSV writing", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to finalize CSV"})
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/csv")
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=all_survey_reports.csv")
	return c.Blob(http.StatusOK, "text/csv", buffer.Bytes())
}

func formatCorrectAnswers(correctAnswers []dto.CorrectAnswerPercentageToShow) string {
	var results []string
	for _, ca := range correctAnswers {
		results = append(results, fmt.Sprintf("%d: %s", ca.QuestionID, ca.Percentage))
	}
	return strings.Join(results, "; ")
}

func formatChoicesPercentage(questionChoices []dto.QuestionReport) string {
	var results []string
	for _, qc := range questionChoices {
		for _, c := range qc.ChoiceReport {
			results = append(results, fmt.Sprintf("%d => %s: %s", qc.QuestionID, c.Text, c.Percentage))
		}
	}
	return strings.Join(results, "; ")
}

func formatDispersionByHour(dispersion []dto.HourDispersionDTO) string {
	var results []string
	for _, hour := range dispersion {
		results = append(results, fmt.Sprintf("%02d:00 - %d", hour.Hour, hour.Count))
	}
	return strings.Join(results, "; ")
}

func formatParticipationReport(participations []dto.ParticipationReport) string {
	var results []string
	for _, p := range participations {
		results = append(results, fmt.Sprintf("%d: %d", p.UserID, p.Count))
	}
	return strings.Join(results, "; ")
}
