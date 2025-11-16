package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"os"
	"path/filepath"
	"savory-ai-server/app/module/qr_code/service"
	"strconv"
)

type qrCodeController struct {
	qrCodeService service.QRCodeService
}

type QRCodeController interface {
	GetRestaurantQRCode(c *fiber.Ctx) error
	GetTableQRCode(c *fiber.Ctx) error
	DownloadRestaurantQRCode(c *fiber.Ctx) error
	DownloadTableQRCode(c *fiber.Ctx) error
}

func NewQRCodeController(service service.QRCodeService) QRCodeController {
	return &qrCodeController{
		qrCodeService: service,
	}
}

func (c *qrCodeController) GetRestaurantQRCode(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("restaurant_qr_%d.png", id)
	qrDir := filepath.Join("storage", "public", "qr_codes", filename)
	if _, err = os.Stat(qrDir); os.IsNotExist(err) {
		return ctx.Status(fiber.StatusNotFound).SendString("image not found")
	}

	return ctx.SendFile(qrDir)
}

func (c *qrCodeController) DownloadRestaurantQRCode(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("restaurant_qr_%d.png", id)
	qrDir := filepath.Join("storage", "public", "qr_codes", filename)

	data, err := os.ReadFile(qrDir)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("image not found")
	}

	contentType := http.DetectContentType(data)
	ctx.Set("Content-Type", contentType)
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	return ctx.SendFile(qrDir)
}

func (c *qrCodeController) GetTableQRCode(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return err
	}
	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("table_qr_%d_%d.png", restaurantID, tableID)
	qrDir := filepath.Join("storage", "public", "qr_codes", filename)
	if _, err = os.Stat(qrDir); os.IsNotExist(err) {
		return ctx.Status(fiber.StatusNotFound).SendString("image not found")
	}

	return ctx.SendFile(qrDir)
}

func (c *qrCodeController) DownloadTableQRCode(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return err
	}
	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("table_qr_%d_%d.png", restaurantID, tableID)
	qrDir := filepath.Join("storage", "public", "qr_codes", filename)

	data, err := os.ReadFile(qrDir)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("image not found")
	}

	contentType := http.DetectContentType(data)
	ctx.Set("Content-Type", contentType)
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	ctx.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	return ctx.SendFile(qrDir)

}
