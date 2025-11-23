package payload

type CreateQuestionReq struct {
	Text         string `json:"text" validate:"required"`
	LanguageCode string `json:"languageCode,omitempty"`
}

type UpdateQuestionReq struct {
	Text         string  `json:"text,omitempty"`
	LanguageCode *string `json:"languageCode,omitempty"`
}
