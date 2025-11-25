// Package controller содержит HTTP-обработчики для модуля чата.
package controller

import (
	"savory-ai-server/app/module/chat/service"
)

// Controller агрегирует все контроллеры модуля чата.
type Controller struct {
	Chat ChatController // Контроллер для работы с чатом
}

// NewControllers создаёт агрегатор контроллеров чата.
// Используется FX для dependency injection.
func NewControllers(chatService service.ChatService) *Controller {
	return &Controller{
		Chat: NewChatController(chatService),
	}
}