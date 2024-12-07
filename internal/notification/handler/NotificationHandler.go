package handler

import (
	"net/http"
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type NotificationHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.INotificationService
	logger  logging.Logger
}

func NewNotificationHandler(conf *config.Config, db db.DbService, logger logging.Logger, notificationService service.INotificationService) *NotificationHandler {
	return &NotificationHandler{conf: conf, db: db, service: notificationService, logger: logger}
}

func (h *NotificationHandler) GetNotifications(c echo.Context) error {

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	var req dto.GetNotificationsRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in  GetNotifications api", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	response, err := h.service.GetNotifications(c.Request().Context(), userID, req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, response)

}

func (h *NotificationHandler) SeenNotification(c echo.Context) error {

	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	id := c.Param("notification_id")

	notifId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Warn(logging.Validation, logging.Api, "validation error in SeenNotification", map[logging.ExtraKey]interface{}{logging.ErrorMessage: err.Error(), logging.UserId: userID})
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid SeenNotification id"})
	}
	err = h.service.Seen(c.Request().Context(), uint(notifId))

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, nil)

}
