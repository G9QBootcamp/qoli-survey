package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type AuthRouter struct {
	conf        *config.Config
	db          db.DbService
	serverGroup *echo.Group
	handler     *handler.UserHandler
	logger      logging.Logger
}

func NewAuthRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger) *AuthRouter {
	return &AuthRouter{conf: conf, db: db, serverGroup: serverGroup, handler: handler.NewHandler(conf, db, logger), logger: logger}
}

func (r *AuthRouter) RegisterRoutes() {
	r.serverGroup.POST("/signup", r.handler.Signup)
	r.serverGroup.POST("/login", r.handler.Login)
}
