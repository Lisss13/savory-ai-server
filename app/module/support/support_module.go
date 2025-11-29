// Package support содержит модуль службы поддержки.
// Позволяет пользователям создавать заявки в поддержку,
// а администраторам - управлять их статусами.
package support

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/support/controller"
	"savory-ai-server/app/module/support/repository"
	"savory-ai-server/app/module/support/service"
)

// SupportRouter маршрутизатор модуля поддержки
type SupportRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

// NewSupportRouter создаёт новый маршрутизатор модуля поддержки
func NewSupportRouter(fiber *fiber.App, controller *controller.Controller) *SupportRouter {
	return &SupportRouter{
		App:        fiber,
		Controller: controller,
	}
}

// SupportModule FX модуль для службы поддержки
var SupportModule = fx.Options(
	fx.Provide(repository.NewSupportRepository),
	fx.Provide(service.NewSupportService),
	fx.Provide(controller.NewSupportController),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewSupportRouter),
)

// RegisterSupportRoutes регистрирует маршруты для пользователей
// Метод: POST /support - создать заявку
// Метод: GET /support/my - мои заявки
// Метод: GET /support/:id - заявка по ID
func (r *SupportRouter) RegisterSupportRoutes(auth fiber.Handler) {
	ctrl := r.Controller.Support
	r.App.Route("/support", func(router fiber.Router) {
		router.Post("/", auth, ctrl.Create)        // Создать заявку
		router.Get("/my", auth, ctrl.GetMyTickets) // Мои заявки
		router.Get("/:id", auth, ctrl.GetByID)     // Заявка по ID
	})
}

// RegisterAdminSupportRoutes регистрирует маршруты для администраторов
// Метод: GET /admin/support - все заявки (с фильтром по статусу)
// Метод: PATCH /admin/support/:id/status - обновить статус
func (r *SupportRouter) RegisterAdminSupportRoutes(auth fiber.Handler, adminAuth fiber.Handler) {
	ctrl := r.Controller.Support
	r.App.Route("/admin/support", func(router fiber.Router) {
		router.Use(auth)      // Сначала проверяем JWT
		router.Use(adminAuth) // Затем проверяем роль admin

		router.Get("/", ctrl.GetAll)               // Все заявки
		router.Patch("/:id/status", ctrl.UpdateStatus) // Обновить статус
	})
}
