package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/organization/payload"
	"savory-ai-server/app/module/organization/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type languageController struct {
	languageService service.LanguageService
}

type LanguageController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	
	// Organization language operations
	GetLanguagesByOrganizationID(c *fiber.Ctx) error
	AddLanguageToOrganization(c *fiber.Ctx) error
	RemoveLanguageFromOrganization(c *fiber.Ctx) error
}

func NewLanguageController(service service.LanguageService) LanguageController {
	return &languageController{
		languageService: service,
	}
}

func (lc *languageController) GetAll(ctx *fiber.Ctx) error {
	languages, err := lc.languageService.GetAllLanguages()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     languages,
		Messages: response.Messages{"Languages retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid language ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	language, err := lc.languageService.GetLanguageByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     language,
		Messages: response.Messages{"Language retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateLanguageReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	language, err := lc.languageService.CreateLanguage(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     language,
		Messages: response.Messages{"Language created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (lc *languageController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid language ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.UpdateLanguageReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	language, err := lc.languageService.UpdateLanguage(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     language,
		Messages: response.Messages{"Language updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid language ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := lc.languageService.DeleteLanguage(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Language deleted successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) GetLanguagesByOrganizationID(ctx *fiber.Ctx) error {
	orgID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid organization ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	languages, err := lc.languageService.GetLanguagesByOrganizationID(uint(orgID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     languages,
		Messages: response.Messages{"Organization languages retrieved successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) AddLanguageToOrganization(ctx *fiber.Ctx) error {
	orgID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid organization ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.AddLanguageToOrgReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := lc.languageService.AddLanguageToOrganization(uint(orgID), req); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Language added to organization successfully"},
		Code:     fiber.StatusOK,
	})
}

func (lc *languageController) RemoveLanguageFromOrganization(ctx *fiber.Ctx) error {
	orgID, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid organization ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.RemoveLanguageFromOrgReq)
	if err := ctx.BodyParser(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid request body"},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	if err := lc.languageService.RemoveLanguageFromOrganization(uint(orgID), req); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Language removed from organization successfully"},
		Code:     fiber.StatusOK,
	})
}