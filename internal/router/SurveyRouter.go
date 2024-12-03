package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type SurveyRouter struct {
	conf    *config.Config
	db      db.DbService
	server  *echo.Echo
	handler *handler.SurveyHandler
	logger  logging.Logger
}

func NewSurveyRouter(conf *config.Config, db db.DbService, server *echo.Echo, logger logging.Logger) *SurveyRouter {
	return &SurveyRouter{conf: conf, db: db, server: server, handler: handler.NewHandler(conf, db, logger), logger: logger}
}

func (r *SurveyRouter) RegisterRoutes() {
	g := r.server.Group("/surveys")
	g.POST("", r.handler.CreateSurvey)
	g.GET("/:survey_id/start", r.handler.StartSurvey)
}
