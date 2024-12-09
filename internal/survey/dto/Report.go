package dto

type CorrectAnswerPercentageToShow struct {
	QuestionID uint   `json:"question_id"`
	Percentage string `json:"percentage"`
}

type QuestionReport struct {
	QuestionID   uint           `json:"question_id"`
	ChoiceReport []ChoiceReport `json:"choice_report"`
}

type ChoiceReport struct {
	ID         uint   `json:"id"`
	Text       string `json:"text"`
	Percentage string `json:"percentage"`
}

type ReportResponse struct {
	SurveyParticipation           string                          `json:"survey_participation"`
	CorrectAnswers                []CorrectAnswerPercentageToShow `json:"correct_answers"`
	MultipleParticipationCount    []ParticipationReport           `json:"multiple_participation_count"`
	SuddenlyFinishedParticipation string                          `json:"suddenly_finished_participation"`
	ChoicesPercentage             []QuestionReport                `json:"choices_percentage"`
	AverageResponseTime           string                          `json:"average_response_time"`
	DispersionResponseByHour      []HourDispersionDTO             `json:"dispersion_response_by_hour"`
}

type ParticipationReport struct {
	UserID uint  `json:"user_id"`
	Count  int64 `json:"count"`
}

type HourDispersionDTO struct {
	Hour  int `json:"hour"`
	Count int `json:"count"`
}
