package service

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/survey/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
)

type IReportService interface {
}
type ReportService struct {
	conf   *config.Config
	repo   repository.ISurveyRepository
	logger logging.Logger
}

func NewReportService(conf *config.Config, repo repository.ISurveyRepository, logger logging.Logger) *ReportService {
	return &ReportService{conf: conf, repo: repo, logger: logger}
}
