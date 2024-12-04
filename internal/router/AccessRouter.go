package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type AccessRouter struct {
	conf    *config.Config
	db      db.DbService
	server  *echo.Echo
	handler *handler.AccessHandler
	logger  logging.Logger
}

func NewAccessRouter(conf *config.Config, db db.DbService, server *echo.Echo, logger logging.Logger) *AccessRouter {
	return &AccessRouter{conf: conf, db: db, server: server, handler: handler.NewAccessHandler(conf, db, logger), logger: logger}
}

func (r *AccessRouter) RegisterRoutes() {
	r.server.GET("/access/permissions", r.handler.GetAllPermissions)
	r.server.GET("/access/assign_role", r.handler.SetRole)
	r.server.GET("/access/get_user_survey_roles", r.handler.GetUserRolesForSomeSurvey)
	r.server.GET("/access/delete_role", r.handler.DeleteUserSurveyRole)
}
