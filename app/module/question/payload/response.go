package payload

import "time"

type LanguageResp struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type QuestionResp struct {
	ID        uint         `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	Text      string       `json:"text"`
	Language  *LanguageResp `json:"language,omitempty"`
}

type QuestionsResp struct {
	Questions []QuestionResp `json:"questions"`
}
