package controller

import "savory-ai-server/app/module/qr_code/service"

type Controller struct {
	QRCode QRCodeController
}

func NewControllers(service service.QRCodeService) *Controller {
	return &Controller{
		QRCode: NewQRCodeController(service),
	}
}
