package handler

import (
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"time"

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
func (h *ReportHandler) WebSocketResults(c echo.Context) error {
	surveyID := c.Param("survey_id")
	userID, ok := c.Request().Context().Value("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, _ := strconv.Atoi(surveyID)

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Failed to upgrade connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}
	defer conn.Close()

	// Channel to signal disconnection
	done := make(chan struct{})

	// Start a goroutine to read messages (to detect disconnection)
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

	// Listen for real-time updates about the survey
	for {
		select {
		case <-done:
			return nil // Exit when the connection is closed
		default:
			// Here you could emit real-time updates on survey responses or status changes
			report, err := h.service.GetReportAggregateService(c.Request().Context(), uint(iSurveyId))
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "Failed to get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return err
			}
			if true {
				// Send details about who voted for what if the survey is not anonymous
				// Example data; replace with actual implementation
				err = conn.WriteJSON(report)
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
					return err
				}
			} else {
				err = conn.WriteJSON(map[string]interface{}{
					"message": "survey is anonymous",
				})
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
					return err
				}
			}
			time.Sleep(2 * time.Second)
		}
	}
}
