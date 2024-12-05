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

	apiGroup := server.Echo.Group("/api")
	apiGroup.Use(middlewares.JWTAuth(conf.JWT.SecretKey))

	userRouter := NewUserRouter(conf, db, apiGroup, logger)
	surveyRouter := NewSurveyRouter(conf, db, apiGroup, logger)
	accessRouter := NewAccessRouter(conf, db, apiGroup, logger)

	authGroup := server.Echo.Group("/auth")
	authRouter := NewAuthRouter(conf, db, authGroup, logger)

	authRouter.RegisterRoutes()
	userRouter.RegisterRoutes()
	surveyRouter.RegisterRoutes()
	accessRouter.RegisterRoutes(db)
	// Additional routers...
}
