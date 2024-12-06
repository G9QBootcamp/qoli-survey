package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/handler"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type QuestionRouter struct {
	conf       *config.Config
	db         db.DbService
	routeGroup *echo.Group
	handler    *handler.QuestionHandler
	logger     logging.Logger
}

func NewQuestionRouter(conf *config.Config, db db.DbService, routeGroup *echo.Group, logger logging.Logger) *QuestionRouter {
	return &QuestionRouter{conf: conf, db: db, routeGroup: routeGroup, handler: handler.NewQuestionHandler(conf, db, logger), logger: logger}
}

func (r *QuestionRouter) RegisterRoutes() {
	g := r.routeGroup.Group("/:survey_id")
	g.GET("/questions", r.handler.GetQuestions)
	g.DELETE("/questions/:question_id", r.handler.DeleteQuestion)
	g.GET("/questions/:question_id", r.handler.GetQuestion)
	g.PATCH("/questions/:question_id", r.handler.UpdateQuestion)

}
