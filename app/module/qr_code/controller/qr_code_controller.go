package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/qr_code/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type qrCodeController struct {
	qrCodeService service.QRCodeService
}

type QRCodeController interface {
	GenerateTableQRCode(c *fiber.Ctx) error
}

func NewQRCodeController(service service.QRCodeService) QRCodeController {
	return &qrCodeController{
		qrCodeService: service,
	}
}

func (c *qrCodeController) GenerateTableQRCode(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(jwt.JWTData)

	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return err
	}

	qrCode, err := c.qrCodeService.GenerateTableQRCode(user.CompanyID, uint(tableID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     qrCode,
		Messages: response.Messages{"QR code generated successfully"},
		Code:     fiber.StatusOK,
	})
}
