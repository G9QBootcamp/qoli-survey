package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type SurveyRouter struct {
	conf        *config.Config
	db          db.DbService
	serverGroup *echo.Group
	handler     *handler.SurveyHandler
	logger      logging.Logger
}

func NewSurveyRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger) *SurveyRouter {
	return &SurveyRouter{conf: conf, db: db, serverGroup: serverGroup, handler: handler.NewSurveyHandler(conf, db, logger), logger: logger}
}

func (r *SurveyRouter) RegisterRoutes() {
	g := r.serverGroup.Group("/surveys")
	g.POST("", r.handler.CreateSurvey)
	g.DELETE("/:survey_id", r.handler.DeleteSurvey, middlewares.CheckPermission("edit_survey", r.db))
	g.GET("/:survey_id", r.handler.GetSurvey, middlewares.CheckPermission("view_survey", r.db))
	g.GET("", r.handler.GetSurveys, middlewares.CheckPermission("view_survey", r.db))
	g.GET("/:survey_id/start", r.handler.StartSurvey, middlewares.CheckPermission("vote", r.db))

	questionRouter := NewQuestionRouter(r.conf, r.db, g, r.logger)
	questionRouter.RegisterRoutes()
}
