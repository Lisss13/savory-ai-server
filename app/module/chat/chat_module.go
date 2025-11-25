// Package chat предоставляет функциональность AI-чата для взаимодействия с посетителями ресторана.
// Поддерживает два типа чатов:
//   - Table Chat: чат для посетителей, сидящих за столиком (через QR-код)
//   - Restaurant Chat: общий чат с рестораном (через AI-бота)
//
// Использует модуль ai_reservation для интеграции с Anthropic Claude.
package chat

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/chat/controller"
	chat_repo "savory-ai-server/app/module/chat/repository"
	"savory-ai-server/app/module/chat/service"
)

// ChatRouter содержит роутер Fiber и контроллеры для обработки HTTP-запросов чата.
type ChatRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

// NewChatRouter создаёт новый экземпляр ChatRouter.
func NewChatRouter(fiber *fiber.App, controller *controller.Controller) *ChatRouter {
	return &ChatRouter{
		App:        fiber,
		Controller: controller,
	}
}

// ChatModule определяет FX-модуль для чата.
// Регистрирует зависимости: репозиторий, сервис чата, контроллер и роутер.
// AI-сервис бронирования предоставляется модулем ai_reservation.
var ChatModule = fx.Options(
	fx.Provide(chat_repo.NewChatRepository),
	fx.Provide(service.NewChatService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewChatRouter),
)

// RegisterChatRoutes регистрирует все маршруты для работы с чатом.
//
// Table Chat (для посетителей за столиком):
//   - GET  /chat/restaurant/:restaurant_id      - получить все чаты ресторана (legacy)
//   - POST /chat/table/session/start            - начать сессию чата для столика
//   - POST /chat/table/session/close/:id        - закрыть сессию чата
//   - POST /chat/table/message/send             - отправить сообщение и получить ответ AI
//   - GET  /chat/table/session/:id/messages     - получить историю сообщений
//   - GET  /chat/table/session/:table_id        - получить сессии для столика
//
// Restaurant Chat (общий чат с рестораном):
//   - POST /chat/restaurant/session/start       - начать сессию чата
//   - POST /chat/restaurant/session/close/:id   - закрыть сессию чата
//   - POST /chat/restaurant/message/send        - отправить сообщение и получить ответ AI
//   - GET  /chat/restaurant/session/:id/messages - получить историю сообщений
//   - GET  /chat/restaurant/sessions/:id        - получить сессии ресторана
func (r *ChatRouter) RegisterChatRoutes() {
	chatController := r.Controller.Chat
	r.App.Route("/chat", func(router fiber.Router) {
		// Table chat endpoints
		router.Get("/restaurant/:restaurant_id", chatController.GetRestaurantChats) // Legacy endpoint
		router.Post("/table/session/start", chatController.StartSessionFromTable)
		router.Post("/table/session/close/:session_id", chatController.CloseSessionFromTable)
		router.Post("/table/message/send", chatController.MessageFromTable)
		router.Get("/table/session/:session_id/messages", chatController.GetMessageFromTableSession)
		router.Get("/table/session/:table_id", chatController.GetSessionsFromTable)

		// Restaurant chat endpoints
		router.Post("/restaurant/session/start", chatController.StartRestaurantSession)
		router.Post("/restaurant/session/close/:session_id", chatController.CloseRestaurantSession)
		router.Post("/restaurant/message/send", chatController.MessageFromRestaurant)
		router.Get("/restaurant/session/:session_id/messages", chatController.GetRestaurantMessagesFromSession)
		router.Get("/restaurant/sessions/:restaurant_id", chatController.GetRestaurantSessions)
	})
}
