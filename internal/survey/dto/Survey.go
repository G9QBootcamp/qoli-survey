package dto

import "time"

type SurveyCreateRequest struct {
	Title              string                  `json:"title" validate:"required"`
	StartTime          time.Time               `json:"start_time" validate:"required"`
	EndTime            time.Time               `json:"end_time" validate:"required"`
	IsSequential       bool                    `json:"is_sequential"`
	AllowReturn        bool                    `json:"allow_return"`
	ParticipationLimit int                     `json:"participation_limit" validate:"required"`
	AnswerTimeLimit    int                     `json:"answer_time_limit" validate:"required"`
	Questions          []QuestionCreateRequest `json:"questions"`
	OwnerID            uint
}

type QuestionCreateRequest struct {
	Text              string                `json:"text" validate:"required"`
	HasMultipleChoice bool                  `json:"has_multiple_choice"`
	MediaUrl          string                `json:"media_url"`
	Choices           []ChoiceCreateRequest `json:"choices"`
	Condition         Condition             `json:"condition"`
}

type ChoiceCreateRequest struct {
	Text      string `json:"text" validate:"required"`
	IsCorrect bool   `json:"is_correct"`
}

type Condition struct {
	QuestionText string `json:"question_text" validate:"required"`
	Answer       string `json:"answer" validate:"required"`
}

type SurveyResponse struct {
	SurveyID           uint       `json:"survey_id"`
	Title              string     `json:"title"`
	StartTime          string     `json:"start_time"`
	EndTime            string     `json:"end_time"`
	IsSequential       bool       `json:"is_sequential"`
	AllowReturn        bool       `json:"allow_return"`
	ParticipationLimit int        `json:"participation_limit"`
	AnswerTimeLimit    int        `json:"answer_time_limit"`
	Questions          []Question `json:"questions"`
}

type Question struct {
	ID                uint     `json:"question_id"`
	Text              string   `json:"text"`
	HasMultipleChoice bool     `json:"has_multiple_choice"`
	MediaUrl          string   `json:"media_url"`
	Choices           []Choice `json:"choices"`
}

type Choice struct {
	ID        uint   `json:"choice_id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}
