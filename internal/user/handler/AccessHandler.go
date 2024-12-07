package handler

import (
	"net/http"
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	notificationService "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type AccessHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IAccessService
	logger  logging.Logger
}

func NewAccessHandler(conf *config.Config, db db.DbService, logger logging.Logger, notificationService notificationService.INotificationService) *AccessHandler {
	return &AccessHandler{conf: conf, db: db, service: service.NewAccessService(conf, repository.NewAccessRepository(db, logger), logger, notificationService), logger: logger}
}

func (h *AccessHandler) SetRole(c echo.Context) error {
	var req dto.SurveyRoleAssignRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	surveyID, err1 := strconv.ParseUint(c.Param("survey_id"), 10, 0)
	userID, err2 := strconv.ParseUint(c.Param("user_id"), 10, 0)

	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	req.SurveyID = uint(surveyID)
	req.UserID = uint(userID)

	res, err := h.service.SetRole(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, res)
}

func (h *AccessHandler) GetUserRolesForSomeSurvey(c echo.Context) error {
	surveyID, err1 := strconv.ParseUint(c.Param("survey_id"), 10, 0)
	userID, err2 := strconv.ParseUint(c.Param("user_id"), 10, 0)

	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	res, err := h.service.GetUserRolesForSomeSurvey(c.Request().Context(), uint(userID), uint(surveyID))
	if err != nil {
		// Handle other errors as internal server errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *AccessHandler) GetAllPermissions(c echo.Context) error {
	res, err := h.service.GetAllPermissions(c.Request().Context())
	if err != nil {
		// Handle other errors as internal server errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (h *AccessHandler) DeleteUserSurveyRole(c echo.Context) error {
	surveyID, err1 := strconv.ParseUint(c.Param("survey_id"), 10, 0)
	userID, err2 := strconv.ParseUint(c.Param("user_id"), 10, 0)
	roleID, err3 := strconv.ParseUint(c.Param("role_id"), 10, 0)

	if err1 != nil || err2 != nil || err3 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	err := h.service.DeleteUserSurveyRole(c.Request().Context(), uint(surveyID), uint(userID), uint(roleID))
	if err != nil {
		// Handle other errors as internal server errors
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)
}
