package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/user/payload"
	"savory-ai-server/app/module/user/service"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
	"strconv"
)

type userController struct {
	userService service.UserService
}

type UserController interface {
	Get(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
}

func NewUserController(service service.UserService) UserController {
	return &userController{
		userService: service,
	}
}

func (uc *userController) Get(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid user ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	user, err := uc.userService.FindUserByID(userID)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusNotFound,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     user,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

func (uc *userController) Update(ctx *fiber.Ctx) error {
	userID, err := strconv.ParseInt(ctx.Params("id"), 10, 64)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{"Invalid user ID"},
			Code:     fiber.StatusBadRequest,
		})
	}

	userData := new(payload.UserUpdateReq)
	if err := ctx.BodyParser(userData); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Validate email if provided
	if userData.Email != "" {
		if err := response.ValidateStruct(userData); err != nil {
			return response.Resp(ctx, response.Response{
				Messages: response.Messages{err.Error()},
				Code:     fiber.StatusBadRequest,
			})
		}
	}

	user, err := uc.userService.UpdateUser(userID, userData)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusNotFound,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     user,
		Messages: response.Messages{"User updated successfully"},
		Code:     fiber.StatusOK,
	})
}

func (uc *userController) Create(ctx *fiber.Ctx) error {
	userData := new(payload.UserCreateReq)
	if err := ctx.BodyParser(userData); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	// Validate request
	if err := response.ValidateStruct(userData); err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	currentUser := ctx.Locals("user").(jwt.JWTData)

	// Create user
	user, err := uc.userService.CreateUser(userData, currentUser.CompanyID)
	if err != nil {
		// Check for duplicate email
		if err.Error() == "email already exists" {
			return response.Resp(ctx, response.Response{
				Messages: response.Messages{err.Error()},
				Code:     fiber.StatusConflict,
			})
		}
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     user,
		Messages: response.Messages{"User created successfully"},
		Code:     fiber.StatusCreated,
	})
}
