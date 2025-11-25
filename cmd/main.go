package main

import (
	fxzerolog "github.com/efectn/fx-zerolog"
	"go.uber.org/fx"
	"savory-ai-server/app/middleware"
	"savory-ai-server/app/module/admin"
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
	"savory-ai-server/app/module/table"
	"savory-ai-server/app/module/user"
	"savory-ai-server/app/router"
	"savory-ai-server/internal/bootstrap"
	"savory-ai-server/internal/bootstrap/database"
	"savory-ai-server/utils/config"
)

func main() {
	fx.New(
		/* provide patterns */
		// config
		fx.Provide(config.NewConfig),
		// logging
		fx.Provide(bootstrap.NewLogger),
		// fiber
		fx.Provide(bootstrap.NewFiber),
		// storage
		fx.Provide(database.NewDatabase),
		// middleware
		fx.Provide(middleware.NewMiddleware),
		// router
		fx.Provide(router.NewRouter),

		// provide modules
		auth.NewAuthModule,
		user.UserModuler,
		menu_category.MenuCategoryModule,
		dish.DishModule,
		file_upload.FileUploadModule,
		qrcode.QRCodeModule,
		table.TableModule,
		question.QuestionModule,
		restaurant.RestaurantModule,
		organization.OrganizationModuler,
		reservation.ReservationModule,
		chat.ChatModule,
		subscription.SubscriptionModule,
		admin.AdminModule,

		// start aplication
		fx.Invoke(bootstrap.Start),

		// define logger
		fx.WithLogger(fxzerolog.Init()),
	).Run()
}
