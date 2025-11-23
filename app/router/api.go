package router

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/middleware"
	"savory-ai-server/app/module/auth"
	"savory-ai-server/app/module/chat"
	"savory-ai-server/app/module/dish"
	"savory-ai-server/app/module/file_upload"
	"savory-ai-server/app/module/menu_category"
	"savory-ai-server/app/module/organization"
	qrcode "savory-ai-server/app/module/qr_code"
	"savory-ai-server/app/module/question"
	"savory-ai-server/app/module/restaurant"
	"savory-ai-server/app/module/subscription"
	"savory-ai-server/app/module/table"
	"savory-ai-server/app/module/user"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/jwt"
)

type Router struct {
	App fiber.Router
	Cfg *config.Config

	AuthRouter         *auth.AuthRouter
	UserRouter         *user.UserRouter
	MenuCategoryRouter *menu_category.MenuCategoryRouter
	DishRouter         *dish.DishRouter
	FileUploadRouter   *file_upload.FileUploadRouter
	QrCodeRouter       *qrcode.QRCodeRouter
	TableRouter        *table.TableRouter
	QuestionRouter     *question.QuestionRouter
	RestaurantRouter   *restaurant.RestaurantRouter
	OrganizationRouter   *organization.OrganizationRouter
	ChatRouter           *chat.ChatRouter
	SubscriptionRouter   *subscription.SubscriptionRouter
}

func NewRouter(
	fiber *fiber.App,
	cfg *config.Config,

	authRouter *auth.AuthRouter,
	userRouter *user.UserRouter,
	menuCategoryRouter *menu_category.MenuCategoryRouter,
	dishRouter *dish.DishRouter,
	fileUploadRouter *file_upload.FileUploadRouter,
	qrCodeRouter *qrcode.QRCodeRouter,
	tableRouter *table.TableRouter,
	questionRouter *question.QuestionRouter,
	restaurantRouter *restaurant.RestaurantRouter,
	organizationRouter *organization.OrganizationRouter,
	chatRouter *chat.ChatRouter,
	subscriptionRouter *subscription.SubscriptionRouter,
) *Router {
	return &Router{
		App:                fiber,
		Cfg:                cfg,
		AuthRouter:         authRouter,
		UserRouter:         userRouter,
		MenuCategoryRouter: menuCategoryRouter,
		DishRouter:         dishRouter,
		FileUploadRouter:   fileUploadRouter,
		QrCodeRouter:       qrCodeRouter,
		TableRouter:        tableRouter,
		QuestionRouter:     questionRouter,
		RestaurantRouter:   restaurantRouter,
		OrganizationRouter:   organizationRouter,
		ChatRouter:           chatRouter,
		SubscriptionRouter:   subscriptionRouter,
	}
}

// Register routes
func (r *Router) Register() {
	// Test Routes
	r.App.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Pong! 👋")
	})

	authRequired := middleware.AuthRequired(r.Cfg)

	r.App.Get("/auth/chek", authRequired, func(c *fiber.Ctx) error {
		user := c.Locals("user").(jwt.JWTData)
		return c.JSON(fiber.Map{
			"success": true,
			"data": fiber.Map{
				"id":         user.ID,
				"email":      user.Email,
				"company_id": user.CompanyID,
			},
		})
	})

	//Register routes of modules
	r.AuthRouter.RegisterAuthRoutes(authRequired)
	r.UserRouter.RegisterUserRouters(authRequired)
	r.MenuCategoryRouter.RegisterMenuCategoryRoutes(authRequired)
	r.DishRouter.RegisterDishRoutes(authRequired)
	r.FileUploadRouter.RegisterFileUploadRoutes()
	r.QrCodeRouter.RegisterQRCodeRoutes(authRequired)
	r.TableRouter.RegisterTableRoutes(authRequired)
	r.QuestionRouter.RegisterQuestionRoutes(authRequired)
	r.RestaurantRouter.RegisterRestaurantRoutes(authRequired)
	r.OrganizationRouter.RegisterOrganizationRoutes(authRequired)
	r.ChatRouter.RegisterChatRoutes()
	r.SubscriptionRouter.RegisterSubscriptionRoutes(authRequired)
}
