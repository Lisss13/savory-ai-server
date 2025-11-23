package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/question/payload"
	"savory-ai-server/app/module/question/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type questionController struct {
	questionService service.QuestionService
}

type QuestionController interface {
	GetAll(c *fiber.Ctx) error
	GetByLanguage(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewQuestionController(service service.QuestionService) QuestionController {
	return &questionController{
		questionService: service,
	}
}

func (c *questionController) GetAll(ctx *fiber.Ctx) error {
	// Get organization ID from a JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	// Check if a language query parameter is provided
	languageCode := ctx.Query("language")

	var questions *payload.QuestionsResp
	var err error

	if languageCode != "" {
		// Get questions by organization ID and language
		questions, err = c.questionService.GetByOrganizationIDAndLanguage(user.CompanyID, languageCode)
	} else {
		// Get all questions for the organization
		questions, err = c.questionService.GetByOrganizationID(user.CompanyID)
	}

	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     questions,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *questionController) GetByLanguage(ctx *fiber.Ctx) error {
	// Get organization ID from a JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	// Get language code from a path parameter
	languageCode := ctx.Params("code")
	if languageCode == "" {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Language code is required"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Get questions by organization ID and language
	questions, err := c.questionService.GetByOrganizationIDAndLanguage(user.CompanyID, languageCode)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     questions,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (c *questionController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateQuestionReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Validate request
	if err := response.ValidateStruct(req); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Get organization ID from JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	// Create the question
	question, err := c.questionService.Create(req, user.CompanyID)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     question,
		Messages: response.Messages{"Question created successfully"},
		Code:     fiber.StatusCreated,
	})
}

func (c *questionController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid question ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	req := new(payload.UpdateQuestionReq)
	if err := ctx.BodyParser(req); err != nil {
		return err
	}

	// Get organization ID from JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	// Update the question
	question, err := c.questionService.Update(uint(id), req, user.CompanyID)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     question,
		Messages: response.Messages{"Question updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (c *questionController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid question ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Get organization ID from JWT token
	user := ctx.Locals("user").(jwt.JWTData)

	if err := c.questionService.Delete(uint(id), user.CompanyID); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data: struct {
			ID uint `json:"id"`
		}{ID: uint(id)},
		Messages: response.Messages{"Question deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
