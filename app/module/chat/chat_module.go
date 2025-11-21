package chat

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/chat/controller"
	chat_repo "savory-ai-server/app/module/chat/repository"
	"savory-ai-server/app/module/chat/service"
)

type ChatRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewChatRouter(fiber *fiber.App, controller *controller.Controller) *ChatRouter {
	return &ChatRouter{
		App:        fiber,
		Controller: controller,
	}
}

var ChatModule = fx.Options(
	fx.Provide(chat_repo.NewChatRepository),
	fx.Provide(service.NewChatService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewChatRouter),
)

func (r *ChatRouter) RegisterChatRoutes(auth fiber.Handler) {
	chatController := r.Controller.Chat
	r.App.Route("/chat", func(router fiber.Router) {
		router.Get("/restaurant/:restaurant_id", chatController.GetRestaurantChats)
		// Message endpoints
		router.Post("/table/session/start", chatController.StartSessionFromTable)
		router.Post("/table/session/close/:session_id", chatController.CloseSessionFromTable)
		router.Post("/table/message/send", chatController.MessageFromTable)
		router.Get("/table/session/:session_id/messages", chatController.GetMessageFromTableSession)
		router.Get("/table/session/:table_id", chatController.GetSessionsFromTable)
	})
}
