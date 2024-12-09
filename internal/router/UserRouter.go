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
	conf        *config.Config
	db          db.DbService
	serverGroup *echo.Group
	handler     *handler.UserHandler
	logger      logging.Logger
}

func NewUserRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger) *UserRouter {
	return &UserRouter{conf: conf, db: db, serverGroup: serverGroup, handler: handler.NewHandler(conf, db, logger), logger: logger}
}

func (r *UserRouter) RegisterRoutes(db db.DbService) {
	r.serverGroup.POST("/restrict/users/:user_id",
		r.handler.RestrictUserSurveys,
		middlewares.CheckPermission("restrict_user", db))
	g := r.serverGroup.Group("/users")
	g.GET("", r.handler.GetUsers)
	g.PATCH("/profile", r.handler.UpdateUserProfile)
	g.GET("/profile", r.handler.GetProfile)

}
