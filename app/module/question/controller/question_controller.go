package controller

import (
	"fmt"
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
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewQuestionController(service service.QuestionService) QuestionController {
	return &questionController{
		questionService: service,
	}
}

func (c *questionController) GetAll(ctx *fiber.Ctx) error {
	questions, err := c.questionService.GetAll()
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

	fmt.Printf("CreateQuestionReq: 0 %v\n", req)
	fmt.Printf("CreateQuestionReq: 1 %v\n", user)

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

func (c *questionController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.questionService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Messages: response.Messages{"Question deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
