package controller

import "savory-ai-server/app/module/menu_category/service"

type Controller struct {
	MenuCategory MenuCategoryController
}

func NewControllers(service service.MenuCategoryService) *Controller {
	return &Controller{
		MenuCategory: NewMenuCategoryController(service),
	}
}
