package controller

import (
	"savory-ai-server/app/module/chat/service"
)

type Controller struct {
	Chat ChatController
}

func NewControllers(chatService service.ChatService) *Controller {
	return &Controller{
		Chat: NewChatController(chatService),
	}
}