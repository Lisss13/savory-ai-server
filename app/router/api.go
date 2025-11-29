package router

import (
	"github.com/gofiber/fiber/v2"
	"savory-ai-server/app/middleware"
	"savory-ai-server/app/module/admin"
	adminMiddleware "savory-ai-server/app/module/admin/middleware"
	"savory-ai-server/app/module/auth"
	"savory-ai-server/app/module/chat"
	"savory-ai-server/app/module/dish"
	"savory-ai-server/app/module/file_upload"
	"savory-ai-server/app/module/menu_category"
	"savory-ai-server/app/module/organization"
	qrcode "savory-ai-server/app/module/qr_code"
	"savory-ai-server/app/module/question"
	"savory-ai-server/app/module/reservation"
	"savory-ai-server/app/module/restaurant"
	"savory-ai-server/app/module/subscription"
	"savory-ai-server/app/module/support"
	"savory-ai-server/app/module/table"
	"savory-ai-server/app/module/user"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/jwt"
	"savory-ai-server/utils/response"
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
	OrganizationRouter *organization.OrganizationRouter
	ReservationRouter  *reservation.ReservationRouter
	ChatRouter         *chat.ChatRouter
	SubscriptionRouter *subscription.SubscriptionRouter
	AdminRouter        *admin.AdminRouter
	SupportRouter      *support.SupportRouter
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
	reservationRouter *reservation.ReservationRouter,
	chatRouter *chat.ChatRouter,
	subscriptionRouter *subscription.SubscriptionRouter,
	adminRouter *admin.AdminRouter,
	supportRouter *support.SupportRouter,
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
		OrganizationRouter: organizationRouter,
		ReservationRouter:  reservationRouter,
		ChatRouter:         chatRouter,
		SubscriptionRouter: subscriptionRouter,
		AdminRouter:        adminRouter,
		SupportRouter:      supportRouter,
	}
}

// Register routes
func (r *Router) Register() {
	// Test Routes
	r.App.Get("/ping", PingHandler)

	authRequired := middleware.AuthRequired(r.Cfg)
	adminRequired := adminMiddleware.AdminRequired()

	r.App.Get("/auth/chek", authRequired, HealthCheckHandler)

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
	r.ReservationRouter.RegisterReservationRoutes(authRequired)
	r.ChatRouter.RegisterChatRoutes()
	r.SubscriptionRouter.RegisterSubscriptionRoutes(authRequired)
	r.AdminRouter.RegisterAdminRoutes(authRequired)
	r.SupportRouter.RegisterSupportRoutes(authRequired)
	r.SupportRouter.RegisterAdminSupportRoutes(authRequired, adminRequired)
}

func PingHandler(c *fiber.Ctx) error {
	return response.RespStatusOk(c, "Pong", "Pong success")
}

func HealthCheckHandler(c *fiber.Ctx) error {
	currentUser := c.Locals("user").(jwt.JWTData)
	res := c.JSON(fiber.Map{
		"id":         currentUser.ID,
		"email":      currentUser.Email,
		"company_id": currentUser.CompanyID,
	})
	return response.RespStatusOk(c, res, "Health check success")
}
