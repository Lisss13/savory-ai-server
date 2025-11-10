package controller

import "savory-ai-server/app/module/file_upload/service"

type Controller struct {
	FileUpload FileUploadController
}

func NewControllers(service service.FileUploadService) *Controller {
	return &Controller{
		FileUpload: NewFileUploadController(service),
	}
}
