package payload

import "time"

type QuestionResp struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
}

type QuestionsResp struct {
	Questions []QuestionResp `json:"questions"`
}