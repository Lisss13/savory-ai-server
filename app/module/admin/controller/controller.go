package controller

import "savory-ai-server/app/module/admin/service"

type Controller struct {
	Admin AdminController
}

func NewControllers(service service.AdminService) *Controller {
	return &Controller{
		Admin: NewAdminController(service),
	}
}
