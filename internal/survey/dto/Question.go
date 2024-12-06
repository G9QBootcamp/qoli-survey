package dto

type GetQuestionsRequest struct {
	SurveyId uint `json:"survey_id" validate:"required"`
}

type Question struct {
	ID                uint     `json:"question_id"`
	Text              string   `json:"text"`
	HasMultipleChoice bool     `json:"has_multiple_choice"`
	MediaUrl          string   `json:"media_url"`
	Choices           []Choice `json:"choices"`
}

type QuestionList []*Question

func (questions QuestionList) GetIds() (ids []uint) {

	for _, v := range questions {
		ids = append(ids, v.ID)
	}

	return ids
}
func (questions QuestionList) ToMap() (mapQuestions map[uint]*Question) {
	mapQuestions = map[uint]*Question{}
	for _, v := range questions {
		mapQuestions[v.ID] = v
	}

	return mapQuestions
}

type Answer string

const NoAnswer Answer = "NoAnswer"

type QuestionsAnswerMap []map[Answer]*Question

type QuestionUpdateRequest struct {
	Text              string                `json:"text" validate:"required"`
	HasMultipleChoice bool                  `json:"has_multiple_choice"`
	MediaUrl          string                `json:"media_url"`
	Choices           []ChoiceUpdateRequest `json:"choices"`
}

type ChoiceUpdateRequest struct {
	Text             string `json:"text" validate:"required"`
	IsCorrect        bool   `json:"is_correct"`
	LinkedQuestionId uint   `json:"linked_question_id"`
}
