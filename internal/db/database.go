package db

import (
	"sync"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"gorm.io/gorm"
)

var once sync.Once
var dbService DbService

type DbService interface {
	Init(cfg *config.Config) error
	Close()
	GetDb() *gorm.DB
}

func New() DbService {
	once.Do(func() {
		dbService = newPostgresDb()

	})
	return dbService
}
