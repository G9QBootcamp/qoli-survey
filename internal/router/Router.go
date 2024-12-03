package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	"github.com/G9QBootcamp/qoli-survey/internal/server"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

func RegisterRoutes(conf *config.Config, db db.DbService, server *server.Server, logger logging.Logger) {
	server.Echo.Use(middlewares.RecoveryErrors(logger))
	server.Echo.Use(middlewares.DefaultStructuredLogger(conf, logger))

	userRouter := NewUserRouter(conf, db, server.Echo, logger)
	surveyRouter := NewSurveyRouter(conf, db, server.Echo, logger)

	userRouter.RegisterRoutes()
	surveyRouter.RegisterRoutes(db)
	// Additional routers...
}
