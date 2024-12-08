package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/server"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

func RegisterRoutes(conf *config.Config, db db.DbService, server *server.Server, logger logging.Logger, notificationService service.INotificationService) {
	server.Echo.Use(middlewares.RecoveryErrors(logger))
	server.Echo.Use(middlewares.DefaultStructuredLogger(conf, logger))

	apiGroup := server.Echo.Group("/api")
	apiGroup.Use(middlewares.JWTAuth(conf.JWT.SecretKey))

	userRouter := NewUserRouter(conf, db, apiGroup, logger)
	surveyRouter := NewSurveyRouter(conf, db, apiGroup, logger, notificationService)
	accessRouter := NewAccessRouter(conf, db, apiGroup, logger, notificationService)
	notificationRouter := NewNotificationRouter(conf, db, apiGroup, logger, notificationService)

	authGroup := server.Echo.Group("/auth")
	authRouter := NewAuthRouter(conf, db, authGroup, logger)

	authRouter.RegisterRoutes()
	userRouter.RegisterRoutes(db)
	surveyRouter.RegisterRoutes()
	accessRouter.RegisterRoutes(db)
	notificationRouter.RegisterRoutes()
	// Additional routers...
}
