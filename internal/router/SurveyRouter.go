package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	middlewares "github.com/G9QBootcamp/qoli-survey/internal/middleware"
	notification "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type SurveyRouter struct {
	conf          *config.Config
	db            db.DbService
	serverGroup   *echo.Group
	handler       *handler.SurveyHandler
	reportHandler *handler.ReportHandler
	logger        logging.Logger
}

func NewSurveyRouter(conf *config.Config, db db.DbService, serverGroup *echo.Group, logger logging.Logger, notificationService notification.INotificationService) *SurveyRouter {
	return &SurveyRouter{conf: conf, db: db, serverGroup: serverGroup,
		handler:       handler.NewSurveyHandler(conf, db, logger, notificationService),
		reportHandler: handler.NewReportHandler(conf, db, logger),
		logger:        logger,
	}
}

func (r *SurveyRouter) RegisterRoutes() {
	g := r.serverGroup.Group("/surveys")
	g.POST("", r.handler.CreateSurvey)
	g.DELETE("/:survey_id", r.handler.DeleteSurvey, middlewares.CheckPermission("edit_survey", r.db))
	g.PATCH("/:survey_id", r.handler.UpdateSurvey, middlewares.CheckPermission("edit_survey", r.db))
	g.GET("/:survey_id", r.handler.GetSurvey, middlewares.CheckPermission("view_survey", r.db))
	g.GET("", r.handler.GetSurveys, middlewares.CheckPermission("view_survey", r.db))
	g.GET("/:survey_id/start", r.handler.StartSurvey, middlewares.CheckPermission("vote", r.db), middlewares.CanUserVoteOnSurvey(r.db))
	g.GET("/:survey_id/reports", r.reportHandler.GetSurveyReport, middlewares.CheckPermission("view_survey_reports", r.db))
	g.POST("/reports-to-csv", r.reportHandler.GenerateAllSurveysReport)
	g.GET("/:survey_id/users/:user_id/votes", r.handler.GetUserVotes)
	g.GET("/:survey_id/visible-vote-users", r.handler.GetVisibleVoteUsers)
	g.GET("/:survey_id/reports/live", r.reportHandler.WebSocketResults, middlewares.CheckPermission("view_survey_reports", r.db))

	g.DELETE("/:survey_id/votes/:vote_id", r.handler.DeleteVote, middlewares.CheckPermission("vote", r.db))
	g.GET("/:survey_id/votes", r.handler.SurveyVotes, middlewares.CheckPermission("view_survey_results", r.db))

	g.POST("/:survey_id/options", r.handler.CreateSurveyOption, middlewares.CheckPermission("edit_survey", r.db))
	g.DELETE("/:survey_id/options/:option_id", r.handler.DeleteSurveyOption, middlewares.CheckPermission("edit_survey", r.db))
	g.PATCH("/:survey_id/options/:option_id", r.handler.UpdateSurveyOption, middlewares.CheckPermission("edit_survey", r.db))
	g.GET("/:survey_id/options", r.handler.GetSurveyOptions, middlewares.CheckPermission("view_survey", r.db))

	questionRouter := NewQuestionRouter(r.conf, r.db, g, r.logger)
	questionRouter.RegisterRoutes()

}
