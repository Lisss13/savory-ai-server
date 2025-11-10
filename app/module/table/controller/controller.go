package controller

import "savory-ai-server/app/module/table/service"

type Controller struct {
	Table TableController
}

func NewControllers(service service.TableService) *Controller {
	return &Controller{
		Table: NewTableController(service),
	}
}
