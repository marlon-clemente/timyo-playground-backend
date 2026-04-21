package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/marlon-clemente/timyo-playground-backend/interface/routes"
	"github.com/marlon-clemente/timyo-playground-backend/packages/config"
	"github.com/marlon-clemente/timyo-playground-backend/packages/database"
	"github.com/marlon-clemente/timyo-playground-backend/packages/logger"
	"github.com/marlon-clemente/timyo-playground-backend/packages/observability"
	"github.com/marlon-clemente/timyo-playground-backend/packages/server"

	_ "github.com/marlon-clemente/timyo-playground-backend/docs"
)

// @title           Timyo Playground API
// @version         1.0
// @description     This is the API server for Timyo Playground.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	ctx := context.Background()

	// 1. Load Configuration
	cfg := config.Load()

	// 2. Initialize Logger
	appLogger := logger.Init(cfg.LogFormat)

	_, err := database.New(cfg.DatabaseDSN, cfg.Development)
	if err != nil {
		appLogger.Error("failed to initialize database", "error", err)
		return
	}

	// 3. Initialize Observability (OTel with OpenObserve sink)
	if cfg.OtelEndpoint != "" {
		shutdown, err := observability.Init(ctx, cfg.ServiceName, cfg.OtelEndpoint, cfg.OtelAuth, cfg.OtelOrg, cfg.OtelInsecure)
		if err != nil {
			appLogger.Error("failed to initialize observability", "error", err)
		} else {
			defer shutdown()
			appLogger.Info("observability initialized", "endpoint", cfg.OtelEndpoint)
		}
	}

	// 4. Initialize App
	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		ErrorHandler:          server.ErrorHandler,
	})

	// 5. Global Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.FrontendBaseUrl,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	app.Use(server.NewLoggingMiddleware(appLogger))
	app.Use(observability.Middleware())

	// 6 . Services
	routes.PrivateRoutes(app, cfg.JWTSecret, cfg.AuthServiceURL)
	routes.PublicRoutes(app)

	// 8. Start Server
	addServer := fmt.Sprintf(":%s", cfg.Port)
	appLogger.Info("running server in  http://localhost" + addServer)
	appLogger.Info("running docs in    http://localhost" + addServer + "/docs")
	if cfg.OtelEndpoint != "" {
		appLogger.Info("running tracing in http://" + cfg.OtelEndpoint)
	} else {
		appLogger.Warn("tracing is disabled")
	}
	log.Fatal(app.Listen(addServer))
}
