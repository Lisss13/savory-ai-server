package qr_code

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"savory-ai-server/app/module/qr_code/controller"
	"savory-ai-server/app/module/qr_code/service"
)

type QRCodeRouter struct {
	App        fiber.Router
	Controller *controller.Controller
}

func NewQRCodeRouter(fiber *fiber.App, controller *controller.Controller) *QRCodeRouter {
	return &QRCodeRouter{
		App:        fiber,
		Controller: controller,
	}
}

var QRCodeModule = fx.Options(
	fx.Provide(service.NewQRCodeService),
	fx.Provide(controller.NewControllers),
	fx.Provide(NewQRCodeRouter),
)

func (r *QRCodeRouter) RegisterQRCodeRoutes(auth fiber.Handler) {
	qrCodeController := r.Controller.QRCode
	r.App.Route("/qrcodes", func(router fiber.Router) {
		router.Get("/restaurant/:restaurant_id", qrCodeController.GetRestaurantQRCode)
		router.Get("/restaurant/:restaurant_id/download", qrCodeController.DownloadRestaurantQRCode)
		router.Get("/restaurant/:restaurant_id/table/:table_id", qrCodeController.GetTableQRCode)
		router.Get("/restaurant/:restaurant_id/table/:table_id/download", qrCodeController.DownloadTableQRCode)
	})
}
