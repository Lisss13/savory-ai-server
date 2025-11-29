package payload

// CreateQuestionReq — запрос на создание нового вопроса.
type CreateQuestionReq struct {
	Text         string `json:"text" validate:"required"`
	LanguageCode string `json:"languageCode,omitempty"`
	// ChatType определяет тип чата: "reservation" (бронирование) или "menu" (меню).
	// По умолчанию "menu".
	ChatType string `json:"chatType,omitempty"`
	// DisplayOrder — порядок отображения вопроса (опционально).
	// Если не указан, вопрос добавляется в конец списка.
	DisplayOrder *int `json:"displayOrder,omitempty"`
}

// UpdateQuestionReq — запрос на обновление вопроса.
type UpdateQuestionReq struct {
	Text         string  `json:"text,omitempty"`
	LanguageCode *string `json:"languageCode,omitempty"`
	// ChatType определяет тип чата: "reservation" или "menu".
	ChatType *string `json:"chatType,omitempty"`
	// DisplayOrder — порядок отображения вопроса.
	DisplayOrder *int `json:"displayOrder,omitempty"`
}

// ReorderQuestionsReq — запрос на изменение порядка отображения вопросов.
// Содержит массив ID вопросов в желаемом порядке.
type ReorderQuestionsReq struct {
	// QuestionIDs — массив ID вопросов в желаемом порядке отображения.
	// Первый элемент получит display_order = 0, второй = 1 и т.д.
	QuestionIDs []uint `json:"questionIds" validate:"required,min=1"`
}
