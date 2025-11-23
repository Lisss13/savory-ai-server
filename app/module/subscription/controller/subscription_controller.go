package controller

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/module/subscription/payload"
	"savory-ai-server/app/module/subscription/service"
	"savory-ai-server/utils/response"
	"strconv"
)

type subscriptionController struct {
	subscriptionService service.SubscriptionService
}

type SubscriptionController interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	GetByOrganizationID(c *fiber.Ctx) error
	GetActiveByOrganizationID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Extend(c *fiber.Ctx) error
	Deactivate(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func NewSubscriptionController(service service.SubscriptionService) SubscriptionController {
	return &subscriptionController{
		subscriptionService: service,
	}
}

// GetAll returns all subscriptions
func (c *subscriptionController) GetAll(ctx *fiber.Ctx) error {
	subscriptions, err := c.subscriptionService.GetAll()
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscriptions,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByID returns a subscription by ID
func (c *subscriptionController) GetByID(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	subscription, err := c.subscriptionService.GetByID(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetByOrganizationID returns all subscriptions for an organization
func (c *subscriptionController) GetByOrganizationID(ctx *fiber.Ctx) error {
	organizationID, err := strconv.ParseUint(ctx.Params("organizationId"), 10, 32)
	if err != nil {
		return err
	}

	subscriptions, err := c.subscriptionService.GetByOrganizationID(uint(organizationID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscriptions,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// GetActiveByOrganizationID returns the active subscription for an organization
func (c *subscriptionController) GetActiveByOrganizationID(ctx *fiber.Ctx) error {
	organizationID, err := strconv.ParseUint(ctx.Params("organizationId"), 10, 32)
	if err != nil {
		return err
	}

	subscription, err := c.subscriptionService.GetActiveByOrganizationID(uint(organizationID))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"success"},
		Code:     fiber.StatusOK,
	})
}

// Create creates a new subscription
func (c *subscriptionController) Create(ctx *fiber.Ctx) error {
	req := new(payload.CreateSubscriptionReq)
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

	subscription, err := c.subscriptionService.Create(req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"Subscription created successfully"},
		Code:     fiber.StatusCreated,
	})
}

// Update updates an existing subscription
func (c *subscriptionController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.UpdateSubscriptionReq)
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

	subscription, err := c.subscriptionService.Update(uint(id), req)
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"Subscription updated successfully"},
		Code:     fiber.StatusOK,
	})
}

// Extend extends an existing subscription
func (c *subscriptionController) Extend(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	req := new(payload.ExtendSubscriptionReq)
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

	subscription, err := c.subscriptionService.Extend(uint(id), req)
	if err != nil {
		return response.Resp(ctx, response.Response{
			Messages: response.Messages{err.Error()},
			Code:     fiber.StatusBadRequest,
		})
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"Subscription extended successfully"},
		Code:     fiber.StatusOK,
	})
}

// Deactivate deactivates a subscription
func (c *subscriptionController) Deactivate(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	subscription, err := c.subscriptionService.Deactivate(uint(id))
	if err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data:     subscription,
		Messages: response.Messages{"Subscription deactivated successfully"},
		Code:     fiber.StatusOK,
	})
}

// Delete deletes a subscription
func (c *subscriptionController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return err
	}

	if err := c.subscriptionService.Delete(uint(id)); err != nil {
		return err
	}

	return response.Resp(ctx, response.Response{
		Data: struct {
			ID uint `json:"id"`
		}{ID: uint(id)},
		Messages: response.Messages{"Subscription deleted successfully"},
		Code:     fiber.StatusOK,
	})
}
