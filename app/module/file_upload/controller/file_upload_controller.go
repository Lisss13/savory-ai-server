package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/file_upload/service"
	"savory-ai-server/utils/response"
)

type fileUploadController struct {
	fileUploadService service.FileUploadService
}

type FileUploadController interface {
	UploadImage(c *fiber.Ctx) error
}

func NewFileUploadController(service service.FileUploadService) FileUploadController {
	return &fileUploadController{
		fileUploadService: service,
	}
}

func (c *fileUploadController) UploadImage(ctx *fiber.Ctx) error {
	// Get the file from the request
	file, err := ctx.FormFile("image")
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Failed to get image from request: " + err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Check if the file is an image
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"File must be an image (jpeg, png, or gif)"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Upload the file
	uploadResp, err := c.fileUploadService.UploadFile(file, "images")
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Failed to upload image: " + err.Error()},
			Code:     fiber.StatusInternalServerError,
		})
	}

	// Return the response
	return response.Resp(ctx, response.Response{
		Data:     uploadResp,
		Messages: response.Messages{"Image uploaded successfully"},
		Code:     fiber.StatusCreated,
	})
}
