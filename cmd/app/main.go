package main

import (
	"log"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
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
}
