package controller

import "savory-ai-server/app/module/subscription/service"

type Controller struct {
	Subscription SubscriptionController
}

func NewControllers(service service.SubscriptionService) *Controller {
	return &Controller{
		Subscription: NewSubscriptionController(service),
	}
}
