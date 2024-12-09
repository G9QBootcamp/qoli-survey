package test

import (
	"context"
	"testing"

	"github.com/G9QBootcamp/qoli-survey/internal/survey/service"
	"github.com/stretchr/testify/assert"
)

func TestGetTotalParticipationPercentage(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	percentage, err := testReportService.GetTotalParticipationPercentage(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetTotalParticipationPercentage: %s", err.Error())
	}

	expectedPercentage := uint(50)
	assert.Equal(t, expectedPercentage, percentage)
}

func TestGetCorrectAnswerPercentage(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	percentage, err := testReportService.GetCorrectAnswerPercentage(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetCorrectAnswerPercentage: %s", err.Error())
	}

	assert.NotNil(t, percentage, "Expected correct answer percentages")
}

func TestSuddenlyFinishedParticipationPercentage(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	percentage, err := testReportService.SuddenlyFinishedParticipationPercentage(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in SuddenlyFinishedParticipationPercentage: %s", err.Error())
	}

	assert.GreaterOrEqual(t, percentage, 0.0, "Suddenly finished participation percentage should be >= 0")
}

func TestGetChoicesByPercentage(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	choices, err := service.GetChoicesByPercentage(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetChoicesByPercentage: %s", err.Error())
	}

	assert.NotNil(t, choices, "Expected choices percentage data")
}

func TestGetMultipleParticipationCount(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	report, err := testReportService.GetMultipleParticipationCount(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetMultipleParticipationCount: %s", err.Error())
	}

	assert.NotNil(t, report, "Expected multiple participation count data")
}

func TestGetAverageResponseTime(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	avgResponseTime, err := testReportService.GetAverageResponseTime(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetAverageResponseTime: %s", err.Error())
	}

	assert.GreaterOrEqual(t, avgResponseTime, 0.0, "Average response time should be >= 0")
}

func TestGetResponseDispersionByHour(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	dispersion, err := testReportService.GetResponseDispersionByHour(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetResponseDispersionByHour: %s", err.Error())
	}
	assert.NotNil(t, dispersion, "Expected response dispersion data by hour")
}

func TestGetSurveyReport(t *testing.T) {
	ctx := context.Background()

	surveyID := uint(1)
	report, err := testReportService.GetSurveyReport(ctx, surveyID)

	if err != nil {
		t.Fatalf("Error in GetSurveyReport: %s", err.Error())
	}

	assert.NotNil(t, report, "Expected survey report data")
}

func TestGetAllSurveys(t *testing.T) {
	ctx := context.Background()

	surveys, err := testReportService.GetAllSurveys(ctx)

	if err != nil {
		t.Fatalf("Error in GetAllSurveys: %s", err.Error())
	}

	assert.Greater(t, len(surveys), 0, "Expected at least one survey")
}

func TestGetAccessibleSurveys(t *testing.T) {
	ctx := context.Background()

	userID := uint(1)
	permission := "view_survey_reports"

	surveys, err := testReportService.GetAccessibleSurveys(ctx, userID, permission)

	if err != nil {
		t.Fatalf("Error in GetAccessibleSurveys: %s", err.Error())
	}

	assert.Greater(t, len(surveys), 0, "Expected at least one accessible survey")
}
