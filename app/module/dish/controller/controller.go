package controller

import (
	"savory-ai-server/app/module/dish/service"
	fileUploadService "savory-ai-server/app/module/file_upload/service"
)

type Controller struct {
	Dish DishController
}

func NewControllers(service service.DishService, fileUploadSvc fileUploadService.FileUploadService) *Controller {
	return &Controller{
		Dish: NewDishController(service, fileUploadSvc),
	}
}
