package controller

import "savory-ai-server/app/module/organization/service"

type Controller struct {
	Organization OrganizationController
	Language     LanguageController
}

func NewControllers(orgService service.OrganizationService, langService service.LanguageService) *Controller {
	return &Controller{
		Organization: NewOrganizationController(orgService),
		Language:     NewLanguageController(langService),
	}
}
