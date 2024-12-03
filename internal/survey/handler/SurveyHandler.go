package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"
	"github.com/gorilla/websocket"

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

func NewHandler(conf *config.Config, db db.DbService, logger logging.Logger) *SurveyHandler {
	return &SurveyHandler{conf: conf, db: db, service: service.New(conf, repository.NewSurveyRepository(db, logger), logger), logger: logger}
}

func (h *SurveyHandler) CreateSurvey(c echo.Context) error {
	var req dto.SurveyCreateRequest

	userID, ok := c.Request().Context().Value("userID").(uint)
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

func (h *SurveyHandler) StartSurvey(c echo.Context) error {

	survey_id := c.Param("survey_id")
	userID, ok := c.Request().Context().Value("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	iSurveyId, _ := strconv.Atoi(survey_id)
	can, canError := h.service.CanUserParticipateToSurvey(c.Request().Context(), userID, uint(iSurveyId))

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

	if canError != nil || !can {
		err := conn.WriteJSON(map[string]string{
			"error": canError.Error(),
		})
		if err != nil {
			h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
			return err
		}
		return nil
	}

	participation, err := h.service.Participate(c.Request().Context(), userID, uint(iSurveyId))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "error in create user participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}
	defer h.service.EndParticipation(c.Request().Context(), participation.ID)

	survey, err := h.service.GetSurvey(c.Request().Context(), uint(iSurveyId))
	if err != nil {
		h.logger.Error(logging.General, logging.Api, "Failed to get survey", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
		return err
	}
	timeLimit := survey.AnswerTimeLimit
	for {
		select {
		case <-done: // Exit when the connection is closed
			err := h.service.EndParticipation(c.Request().Context(), participation.ID)
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in ending user survey participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return err
			}
			return nil

		default:

			if timeLimit <= 0 {
				err := h.service.EndParticipation(c.Request().Context(), participation.ID)
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in ending user survey participation", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
					return err
				}
				err = conn.WriteJSON(map[string]string{
					"message": "survey time limit reached",
				})
				if err != nil {
					h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
					return err
				}
				return nil
			}
			err := conn.WriteJSON(map[string]int{
				"remaining_time": timeLimit,
			})
			if err != nil {
				h.logger.Error(logging.General, logging.Api, "error in writing to websocket connection", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error()})
				return err
			}
			timeLimit--
			time.Sleep(1 * time.Second)
		}
	}

}
