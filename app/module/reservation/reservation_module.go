// Package reservation предоставляет функциональность для управления бронированием столиков в ресторане.
// Модуль поддерживает как публичные эндпоинты для анонимных пользователей (через AI-чат),
// так и защищённые эндпоинты для персонала ресторана.
package reservation

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/reservation/controller"
	"savory-ai-server/app/module/reservation/repository"
	"savory-ai-server/app/module/reservation/service"
)

// ReservationRouter содержит роутер Fiber и контроллеры для обработки HTTP-запросов бронирования.
type ReservationRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

// NewReservationRouter создаёт новый экземпляр ReservationRouter.
// Используется FX для dependency injection.
func NewReservationRouter(fiber *fiber.App, controller *controller.Controller) *ReservationRouter {
	return &ReservationRouter{
		App:        fiber,
		Controller: controller,
	}
}

// ReservationModule определяет FX-модуль для бронирования столиков.
// Регистрирует все зависимости: репозиторий, сервис, контроллер и роутер.
var ReservationModule = fx.Options(
	fx.Provide(repository.NewReservationRepository),
	fx.Provide(service.NewReservationService),
	fx.Provide(controller.NewReservationController),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewReservationRouter),
)

// RegisterReservationRoutes регистрирует все маршруты для работы с бронированием.
//
// Публичные маршруты (без авторизации - для анонимных посетителей через AI-чат):
//   - GET  /reservations/available/:restaurant_id - получить доступные слоты
//   - GET  /reservations/my?phone=...             - получить брони по номеру телефона
//   - POST /reservations/                         - создать бронирование
//   - POST /reservations/:id/cancel/public?phone= - отменить бронь по телефону
//
// Защищённые маршруты (требуют авторизации - для персонала ресторана):
//   - GET    /reservations/                       - получить все бронирования
//   - GET    /reservations/:id                    - получить бронь по ID
//   - GET    /reservations/restaurant/:id         - получить брони ресторана
//   - PATCH  /reservations/:id                    - обновить бронирование
//   - POST   /reservations/:id/cancel             - отменить бронь (админ)
//   - DELETE /reservations/:id                    - удалить бронирование
func (r *ReservationRouter) RegisterReservationRoutes(auth fiber.Handler) {
	ctrl := r.Controller.Reservation

	r.App.Route("/reservations", func(router fiber.Router) {
		// Public routes (for anonymous customers via chat)
		router.Get("/available/:restaurant_id", ctrl.GetAvailableSlots)
		router.Get("/my", ctrl.GetByPhone)                    // Get reservations by phone (query: ?phone=...)
		router.Post("/", ctrl.Create)                         // Create reservation (no auth required)
		router.Post("/:id/cancel/public", ctrl.CancelByPhone) // Cancel by phone (query: ?phone=...)

		// Protected routes (for restaurant staff)
		router.Get("/", auth, ctrl.GetAll)
		router.Get("/:id", auth, ctrl.GetByID)
		router.Get("/restaurant/:restaurant_id", auth, ctrl.GetByRestaurantID)
		router.Patch("/:id", auth, ctrl.Update)
		router.Post("/:id/cancel", auth, ctrl.Cancel)
		router.Delete("/:id", auth, ctrl.Delete)
	})
}
