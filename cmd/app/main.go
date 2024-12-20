package main

import (
	"fmt"
	"log"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/db/seeds"
	"github.com/G9QBootcamp/qoli-survey/internal/notification/repository"
	notification "github.com/G9QBootcamp/qoli-survey/internal/notification/service"
	"github.com/G9QBootcamp/qoli-survey/internal/router"
	"github.com/G9QBootcamp/qoli-survey/internal/server"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

func main() {

	conf, err := config.Load()
	if err != nil {
		log.Fatalf("load config error: %v", err)
	}

	logger := logging.NewLogger(conf)

	dbService := db.New()
	err = dbService.Init(conf)
	if err != nil {
		logger.Fatal(logging.Database, logging.Startup, "error in initializing database", map[logging.ExtraKey]interface{}{logging.Service: "Database", logging.ErrorMessage: err.Error()})
	}
	defer dbService.Close()
	db, err := dbService.GetDb().DB()
	if err != nil {
		logger.Fatal(logging.Database, logging.Startup, "error in initializing database", map[logging.ExtraKey]interface{}{logging.Service: "Database", logging.ErrorMessage: err.Error()})
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal(logging.Database, logging.Startup, "error in initializing database", map[logging.ExtraKey]interface{}{logging.Service: "Database", logging.ErrorMessage: err.Error()})
	}

	seeder := seeds.NewSeeder(dbService, logger)
	seeder.RunSeeders()

	notificationService := notification.New(conf, repository.NewNotificationRepository(dbService, logger), logger)
	s := server.NewHttpServer()
	router.RegisterRoutes(conf, dbService, s, logger, notificationService)

	s.Start(fmt.Sprintf("%s:%d", conf.HTTP.Host, conf.HTTP.Port))

}
