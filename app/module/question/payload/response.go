package payload

import "time"

// LanguageResp — информация о языке в ответе API.
type LanguageResp struct {
	ID   uint   `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// QuestionResp — информация о вопросе в ответе API.
type QuestionResp struct {
	ID        uint          `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	Text      string        `json:"text"`
	Language  *LanguageResp `json:"language,omitempty"`
	// ChatType — тип чата: "reservation" (бронирование) или "menu" (меню).
	ChatType string `json:"chat_type"`
}

// QuestionsResp — список вопросов в ответе API.
type QuestionsResp struct {
	Questions []QuestionResp `json:"questions"`
}
