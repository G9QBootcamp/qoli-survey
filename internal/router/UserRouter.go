package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	"github.com/G9QBootcamp/qoli-survey/internal/user/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type UserRouter struct {
	conf    *config.Config
	db      db.DbService
	server  *echo.Echo
	handler *handler.UserHandler
	logger  logging.Logger
}

func NewUserRouter(conf *config.Config, db db.DbService, server *echo.Echo, logger logging.Logger) *UserRouter {
	return &UserRouter{conf: conf, db: db, server: server, handler: handler.NewHandler(conf, db, logger), logger: logger}
}

func (r *UserRouter) RegisterRoutes(db db.DbService) {
	r.server.GET("/users", r.handler.GetUsers)
	r.server.POST("/signup", r.handler.Signup)
	r.server.POST("/restrict/user/:user_id",
		r.handler.RestrictUserSurveys,
		middlewares.CheckPermission("restrict_user", db))
}
