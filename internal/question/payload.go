package question

type BaseQuestionEvent struct {
	Action string `json:"action"`
}

type QuestionCreatedEvent struct {
	Action           string `json:"action"`
	ProductVariantID uint   `json:"product_variant_id"`
	UserID           uint   `json:"user_id"`
	QuestionText     string `json:"question_text"`
}

type QuestionGetEvent struct {
	ProductVariantID uint `json:"product_variant_id"`
}

type QuestionAnsweredEvent struct {
	Action     string `json:"action"`
	QuestionID uint   `json:"question_id"`
	AnswerText string `json:"answer_text"`
}

type QuestionDeletedEvent struct {
	Action     string `json:"action"`
	QuestionID uint   `json:"question_id"`
}