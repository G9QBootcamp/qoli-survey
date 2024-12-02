package test

import (
	// "magical-crwler/config"

	"os"
	"testing"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/gommon/log"
)

var testDbService db.DbService
var testUserRepo repository.IUserRepository
var testUserService service.IUserService

func TestMain(m *testing.M) {

	testDbService = db.New()
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config:%s ", err.Error())
	}

	logger := logging.NewLogger(cfg)

	err = testDbService.Init(cfg)
	if err != nil {
		logger.Fatalf("Failed to connect to database:%v ", err.Error())
	}
	defer testDbService.Close()
	testUserRepo = repository.NewUserRepository(testDbService, logger)
	testUserService = service.New(cfg, testUserRepo, logger)
	code := m.Run()

	os.Exit(code)

}
