package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/dto"

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

func NewHandler(conf *config.Config, db db.DbService, logger logging.Logger) *SurveyHandler {
	return &SurveyHandler{conf: conf, db: db, service: service.New(conf, repository.NewSurveyRepository(db, logger))}
}

func (h *SurveyHandler) CreateSurvey(c echo.Context) error {
	var req dto.SurveyCreateRequest

	userID, ok := c.Request().Context().Value("userID").(uint)
	if !ok || userID == 0 {
		//return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}
	userID = 1
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	req.OwnerID = userID
	survey, err := h.service.CreateSurvey(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, survey)
}
