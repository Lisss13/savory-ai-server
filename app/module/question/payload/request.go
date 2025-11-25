package payload

// CreateQuestionReq — запрос на создание нового вопроса.
type CreateQuestionReq struct {
	Text         string `json:"text" validate:"required"`
	LanguageCode string `json:"languageCode,omitempty"`
	// ChatType определяет тип чата: "reservation" (бронирование) или "menu" (меню).
	// По умолчанию "menu".
	ChatType string `json:"chatType,omitempty"`
}

// UpdateQuestionReq — запрос на обновление вопроса.
type UpdateQuestionReq struct {
	Text         string  `json:"text,omitempty"`
	LanguageCode *string `json:"languageCode,omitempty"`
	// ChatType определяет тип чата: "reservation" или "menu".
	ChatType *string `json:"chatType,omitempty"`
}
