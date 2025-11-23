package bootstrap

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net"
	"os"
	"runtime"
	"savory-ai-server/app/middleware"
	"savory-ai-server/app/router"
	"savory-ai-server/app/storage"
	"savory-ai-server/internal/bootstrap/database"
	"savory-ai-server/utils/config"
	"savory-ai-server/utils/response"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

// initialize the webserver
func NewFiber(cfg *config.Config) *fiber.App {
	// setup
	app := fiber.New(fiber.Config{
		ServerHeader:          cfg.App.Name,
		AppName:               cfg.App.Name,
		Prefork:               cfg.App.Prefork,
		ErrorHandler:          response.ErrorHandler,
		IdleTimeout:           cfg.App.IdleTimeout * time.Second,
		EnablePrintRoutes:     cfg.App.PrintRoutes,
		DisableStartupMessage: true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		MaxAge:       3600,
	}))

	// pass production config to check it
	response.IsProduction = cfg.App.Production

	return app
}

// function to start webserver
func Start(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	fiber *fiber.App,
	router *router.Router,
	middlewares *middleware.Middleware,
	database *database.Database,
	log zerolog.Logger,
) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {

				// Register middlewares & routes
				middlewares.Register()
				router.Register()

				// Information message
				log.Info().Msg(fiber.Config().AppName + " is running at the moment!")

				// Debug informations
				if !cfg.App.Production {
					prefork := "Enabled"
					procs := runtime.GOMAXPROCS(0)
					if !cfg.App.Prefork {
						procs = 1
						prefork = "Disabled"
					}

					log.Debug().Msgf("Version: %s", "-")
					log.Debug().Msgf("Host: %s", cfg.App.Host)
					log.Debug().Msgf("Port: %s", cfg.App.Port)
					log.Debug().Msgf("Prefork: %s", prefork)
					log.Debug().Msgf("Handlers: %d", fiber.HandlersCount())
					log.Debug().Msgf("Processes: %d", procs)
					log.Debug().Msgf("PID: %d", os.Getpid())
				}

				// Listen the app (with TLS Support)
				if cfg.App.TLS.Enable {
					log.Debug().Msg("TLS support was enabled.")

					if err := fiber.ListenTLS(cfg.App.Port, cfg.App.TLS.CertFile, cfg.App.TLS.KeyFile); err != nil {
						log.Error().Err(err).Msg("An unknown error occurred when to run server!")
					}
				}

				addr := net.JoinHostPort(cfg.App.Host, strings.TrimPrefix(cfg.App.Port, ":"))

				go func() {
					if err := fiber.Listen(addr); err != nil {
						log.Error().Err(err).Msg("An unknown error occurred when to run server 1!")
					}
				}()

				database.ConnectDatabase()

				migrate := flag.Bool("migrate", false, "migrate the storage")
				seeder := flag.Bool("seed", false, "seed the storage")
				flag.Parse()

				// read flag -migrate to migrate the storage
				if *migrate {
					database.MigrateModels()
				}
				// read flag -seed to seed the storage
				if *seeder {
					database.SeedModels(
						&storage.LanguageSeeder{DB: database.DB},
					)
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msg("Shutting down the app...")
				if err := fiber.Shutdown(); err != nil {
					log.Panic().Err(err).Msg("")
				}

				log.Info().Msg("Running cleanup tasks...")
				log.Info().Msg("1- Shutdown the storage")
				database.ShutdownDatabase()
				log.Info().Msgf("%s was successful shutdown.", cfg.App.Name)
				log.Info().Msg("\u001b[96msee you again👋\u001b[0m")

				return nil
			},
		},
	)
}
