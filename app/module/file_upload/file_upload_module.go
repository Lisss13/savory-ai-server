package file_upload

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/file_upload/controller"
	"savory-ai-server/app/module/file_upload/service"
)

type FileUploadRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewFileUploadRouter(fiber *fiber.App, controller *controller.Controller) *FileUploadRouter {
	return &FileUploadRouter{
		App:        fiber,
		Controller: controller,
	}
}

var FileUploadModule = fx.Options(
	fx.Provide(service.NewFileUploadService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewFileUploadRouter),
)

func (r *FileUploadRouter) RegisterFileUploadRoutes() {
	fileUploadController := r.Controller.FileUpload
	r.App.Route("/uploads", func(router fiber.Router) {
		router.Post("/images", fileUploadController.UploadImage)
	})
}
