package test

import (
	// "magical-crwler/config"

	"os"
	"testing"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	surveyRepository "github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	surveyService "github.com/G9QBootcamp/qoli-survey/internal/survey/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/gommon/log"
)

var testDbService db.DbService
var testUserRepo repository.IUserRepository
var testUserService service.IUserService

var testReportRepo surveyRepository.IReportRepository
var testReportService surveyService.IReportService

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

	testReportRepo = surveyRepository.NewReportRepository(testDbService, logger)
	testReportService = surveyService.NewReportService(cfg, testReportRepo, logger)

	code := m.Run()

	os.Exit(code)

}
