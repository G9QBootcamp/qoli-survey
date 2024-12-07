package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/handler"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type NotificationRouter struct {
	conf        *config.Config
	db          db.DbService
	serverGroup *echo.Group
	handler     *handler.NotificationHandler
	logger      logging.Logger
}

func NewNotificationRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger, notificationService service.INotificationService) *NotificationRouter {
	return &NotificationRouter{conf: conf, db: db, serverGroup: serverGroup,
		handler: handler.NewNotificationHandler(conf, db, logger, notificationService),
		logger:  logger,
	}
}

func (r *NotificationRouter) RegisterRoutes() {
	g := r.serverGroup.Group("/notifications")
	g.GET("", r.handler.GetNotifications)
	g.POST("/:notification_id/seen", r.handler.SeenNotification)
}
