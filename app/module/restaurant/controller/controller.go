package controller

import "savory-ai-server/app/module/restaurant/service"

type Controller struct {
	Restaurant RestaurantController
}

func NewControllers(service service.RestaurantService) *Controller {
	return &Controller{
		Restaurant: NewRestaurantController(service),
	}
}
