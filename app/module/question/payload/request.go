package payload

type CreateQuestionReq struct {
	Text string `json:"text" validate:"required"`
}