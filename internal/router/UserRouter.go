package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
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

func (r *UserRouter) RegisterRoutes() {
	r.server.GET("/users", r.handler.GetUsers)
}
