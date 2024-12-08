package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type AccessRouter struct {
	conf        *config.Config
	db          db.DbService
	serverGroup *echo.Group
	handler     *handler.AccessHandler
	logger      logging.Logger
}

func NewAccessRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger, notificationService service.INotificationService) *AccessRouter {
	return &AccessRouter{conf: conf, db: db, serverGroup: serverGroup, handler: handler.NewAccessHandler(conf, db, logger, notificationService), logger: logger}
}

func (r *AccessRouter) RegisterRoutes(db db.DbService) {
	r.serverGroup.GET("/access/permissions", r.handler.GetAllPermissions)
	r.serverGroup.POST("/access/survey/:survey_id/user/:user_id/assign-role",
		r.handler.SetRole,
		middlewares.CheckPermission("assign_and_remove_survey_roles", db),
	)
	r.serverGroup.GET("/access/survey/:survey_id/user/:user_id/roles",
		r.handler.GetUserRolesForSomeSurvey,
		middlewares.CheckPermission("assign_and_remove_survey_roles", db),
	)
	r.serverGroup.DELETE("/access/survey/:survey_id/user/:user_id/role/:role_id",
		r.handler.DeleteUserSurveyRole,
		middlewares.CheckPermission("assign_and_remove_survey_roles", db),
	)
}
