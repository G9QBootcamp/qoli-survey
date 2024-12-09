package dto

import (
	"time"

	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
)

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

type SurveyUpdateRequest struct {
	Title              string    `json:"title" validate:"required"`
	StartTime          time.Time `json:"start_time" validate:"required"`
	EndTime            time.Time `json:"end_time" validate:"required"`
	IsSequential       bool      `json:"is_sequential"`
	AllowReturn        bool      `json:"allow_return"`
	ParticipationLimit int       `json:"participation_limit" validate:"required"`
	AnswerTimeLimit    int       `json:"answer_time_limit" validate:"required"`
}

type SurveysGetRequest struct {
	Page   int    `query:"page" validate:"numeric"`
	UserId int    `query:"page" validate:"numeric"`
	Title  string `query:"title"`
}

type SurveyOptionCreateRequest struct {
	Name  string `json:"name" validate:"required"`
	Value string `json:"value" validate:"required"`
}
type SurveyOptionResponse struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
type SurveyOptionsGetRequest struct {
	SurveyId uint   `json:"survey_id" validate:"numeric"`
	Name     string `json:"name" validate:"required"`
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
	SurveyID           uint                   `json:"survey_id"`
	UserId             uint                   `json:"user_id"`
	Title              string                 `json:"title"`
	StartTime          string                 `json:"start_time"`
	EndTime            string                 `json:"end_time"`
	IsSequential       bool                   `json:"is_sequential"`
	AllowReturn        bool                   `json:"allow_return"`
	ParticipationLimit int                    `json:"participation_limit"`
	AnswerTimeLimit    int                    `json:"answer_time_limit"`
	Options            []SurveyOptionResponse `json:"options"`
}

type Choice struct {
	ID               uint   `json:"choice_id"`
	Text             string `json:"text"`
	IsCorrect        bool   `json:"is_correct"`
	LinkedQuestionID uint   `json:"linked_question_id"`
}

type UserSurveyParticipationResponse struct {
	ID          uint      `json:"id"`
	UserId      uint      `json:"user_id"`
	SurveyID    uint      `json:"survey_id"`
	StartAt     time.Time `json:"start_at"`
	EndAt       time.Time `json:"end_at"`
	CommittedAt time.Time `json:"committed_at"`
}

type OperationType string

const CommitOperation OperationType = "commit"
const BackOperation OperationType = "back"

type VoteRequest struct {
	Operation  OperationType `json:"operation"`
	QuestionId uint          `json:"question_id" validate:"numeric"`
	Answer     string        `json:"answer"`
}
type VoteResponse struct {
	Question *Question `json:"question"`
	Message  string    `json:"message"`
}

type GetVoteResponse struct {
	ID         uint             `json:"id"`
	VoterID    uint             `json:"voter_id"`
	QuestionID uint             `json:"question_id"`
	Answer     string           `json:"answer"`
	IsCorrect  bool             `json:"is_correct"`
	Voter      dto.UserResponse `json:"voter"`
	CreatedAt  time.Time        `json:"created_at"`
}
