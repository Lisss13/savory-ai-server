package controller

import "savory-ai-server/app/module/user/service"

type Controller struct {
	User UserController
}

func NewControllers(service service.UserService) *Controller {
	return &Controller{
		User: NewUserController(service),
	}
}
