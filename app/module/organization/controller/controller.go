package controller

import "savory-ai-server/app/module/organization/service"

type Controller struct {
	Organization OrganizationController
}

func NewControllers(service service.OrganizationService) *Controller {
	return &Controller{
		Organization: NewOrganizationController(service),
	}
}
