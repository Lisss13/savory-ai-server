package controller

import "savory-ai-server/app/module/question/service"

type Controller struct {
	Question QuestionController
}

func NewControllers(service service.QuestionService) *Controller {
	return &Controller{
		Question: NewQuestionController(service),
	}
}
