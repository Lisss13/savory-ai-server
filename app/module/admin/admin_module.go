package admin

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/admin/controller"
	adminMiddleware "savory-ai-server/app/module/admin/middleware"
	admin_repo "savory-ai-server/app/module/admin/repository"
	"savory-ai-server/app/module/admin/service"
)

type AdminRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewAdminRouter(fiber *fiber.App, controller *controller.Controller) *AdminRouter {
	return &AdminRouter{
		App:        fiber,
		Controller: controller,
	}
}

var AdminModule = fx.Options(
	fx.Provide(admin_repo.NewAdminRepository),
	fx.Provide(service.NewAdminService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewAdminRouter),
)

func (r *AdminRouter) RegisterAdminRoutes(auth fiber.Handler) {
	adminController := r.Controller.Admin
	adminRequired := adminMiddleware.AdminRequired()

	// Все роуты админ-панели требуют авторизации + роль admin
	r.App.Route("/admin", func(router fiber.Router) {
		// Применяем оба middleware: auth и adminRequired
		router.Use(auth)
		router.Use(adminRequired)

		// Статистика (дашборд)
		router.Get("/stats", adminController.GetStats)

		// Управление пользователями
		router.Get("/users", adminController.GetAllUsers)
		router.Get("/users/:id", adminController.GetUserByID)
		router.Patch("/users/:id/status", adminController.UpdateUserStatus)
		router.Patch("/users/:id/role", adminController.UpdateUserRole)
		router.Delete("/users/:id", adminController.DeleteUser)

		// Управление организациями
		router.Get("/organizations", adminController.GetAllOrganizations)
		router.Get("/organizations/:id", adminController.GetOrganizationByID)
		router.Delete("/organizations/:id", adminController.DeleteOrganization)

		// Модерация контента (блюда)
		router.Get("/dishes", adminController.GetAllDishes)
		router.Delete("/dishes/:id", adminController.DeleteDish)

		// Логи действий администраторов
		router.Get("/logs", adminController.GetAllLogs)
		router.Get("/logs/me", adminController.GetMyLogs)
	})
}
