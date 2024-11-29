package main

import (
	"fmt"
	"log"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/router"
	"github.com/G9QBootcamp/qoli-survey/internal/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	dbService := db.New()
	dbService.Init(conf)
	defer dbService.Close()
	db, err := dbService.GetDb().DB()
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	s := server.NewHttpServer()
	registerRouters(conf, dbService, s)

	s.Start(fmt.Sprintf("%s:%d", conf.HTTP.Host, conf.HTTP.Port))

}

func registerRouters(conf *config.Config, db db.DbService, server *server.Server) {
	userRouter := router.NewUserRouter(conf, db, server.Echo)
	userRouter.RegisterRoutes()

	//...
}
