package dto

type CorrectAnswerPercentageToShow struct {
	QuestionID       uint    `json:"question_id"`
	HasCorrectAnswer bool    `json:"has_correct_answer"`
	Percentage       float64 `json:"percentage"`
}

type QuestionReport struct {
	QuestionID   uint           `json:"question_id"`
	ChoiceReport []ChoiceReport `json:"choice_report"`
}

type ChoiceReport struct {
	ID         uint    `json:"id"`
	Text       string  `json:"text"`
	Percentage float64 `json:"percentage"`
}

type ReportResponse struct {
	SurveyParticipation           string                          `json:"survey_participation"`
	CorrectAnswers                []CorrectAnswerPercentageToShow `json:"correct_answers"`
	SuddenlyFinishedParticipation string                          `json:"suddenly_finished_participation"`
	ChoicesPercentage             []QuestionReport                `json:"choices_percentage"`
}
