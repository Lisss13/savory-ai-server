package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/qr_code/service"
)

// qrBasePath базовый путь к директории с QR-кодами.
var qrBasePath string

func init() {
	// Получаем абсолютный путь к директории с QR-кодами
	wd, _ := os.Getwd()
	qrBasePath = filepath.Join(wd, "storage", "public", "qr_codes")
}

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

// GetRestaurantQRCode отдаёт QR-код ресторана как статический файл.
// Метод: GET /qrcodes/restaurant/:restaurant_id
func (c *qrCodeController) GetRestaurantQRCode(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid restaurant_id")
	}

	filename := fmt.Sprintf("restaurant_qr_%d.png", id)
	filePath := filepath.Join(qrBasePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	// CORS и заголовки для статического файла
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")

	return ctx.Send(data)
}

// DownloadRestaurantQRCode скачивание QR-кода ресторана.
// Метод: GET /qrcodes/restaurant/:restaurant_id/download
func (c *qrCodeController) DownloadRestaurantQRCode(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid restaurant_id")
	}

	filename := fmt.Sprintf("restaurant_qr_%d.png", id)
	filePath := filepath.Join(qrBasePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	// CORS и заголовки для скачивания
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Content-Type", "image/png")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return ctx.Send(data)
}

// GetTableQRCode отдаёт QR-код столика как статический файл.
// Метод: GET /qrcodes/restaurant/:restaurant_id/table/:table_id
func (c *qrCodeController) GetTableQRCode(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid restaurant_id")
	}
	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid table_id")
	}

	filename := fmt.Sprintf("table_qr_%d_%d.png", restaurantID, tableID)
	filePath := filepath.Join(qrBasePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	// CORS и заголовки для статического файла
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")

	return ctx.Send(data)
}

// DownloadTableQRCode скачивание QR-кода столика.
// Метод: GET /qrcodes/restaurant/:restaurant_id/table/:table_id/download
func (c *qrCodeController) DownloadTableQRCode(ctx *fiber.Ctx) error {
	restaurantID, err := strconv.ParseUint(ctx.Params("restaurant_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid restaurant_id")
	}
	tableID, err := strconv.ParseUint(ctx.Params("table_id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid table_id")
	}

	filename := fmt.Sprintf("table_qr_%d_%d.png", restaurantID, tableID)
	filePath := filepath.Join(qrBasePath, filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("QR code not found")
	}

	// CORS и заголовки для скачивания
	ctx.Set("Access-Control-Allow-Origin", "*")
	ctx.Set("Content-Type", "image/png")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return ctx.Send(data)
}
