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
	Percentage float64 `json:"percentage"`
}
