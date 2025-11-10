package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/organization/payload"
	"savory-ai-server/app/module/organization/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type organizationController struct {
	organizationService service.OrganizationService
}

type OrganizationController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	AddUser(c *fiber.Ctx) error
	RemoveUser(c *fiber.Ctx) error
}

func NewOrganizationController(service service.OrganizationService) OrganizationController {
	return &organizationController{
		organizationService: service,
	}
}

func (oc *organizationController) GetAll(ctx *fiber.Ctx) error {
	organizations, err := oc.organizationService.GetAllOrganizations()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     organizations,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (oc *organizationController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	organization, err := oc.organizationService.GetOrganizationByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     organization,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (oc *organizationController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateOrganizationReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	organization, err := oc.organizationService.UpdateOrganization(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     organization,
		Messages: response.Messages{"Organization updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (oc *organizationController) AddUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.AddUserToOrgReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err = oc.organizationService.AddUserToOrganization(uint(id), req.UserID); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"User added to organization successfully"},
		Code:     fiber.StatusOK,
	})
}

func (oc *organizationController) RemoveUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.RemoveUserFromOrgReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := oc.organizationService.RemoveUserFromOrganization(uint(id), req); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"User removed from organization successfully"},
		Code:     fiber.StatusOK,
	})
}
