package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/server"
)

func RegisterRoutes(conf *config.Config, db db.DbService, server *server.Server) {
	userRouter := NewUserRouter(conf, db, server.Echo)
	userRouter.RegisterRoutes()
	// Additional routers...
}
