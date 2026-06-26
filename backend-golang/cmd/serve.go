package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/spf13/cobra"
	"react-example/backend-golang/config"
	"react-example/backend-golang/internal/handlers"
	customLogger "react-example/backend-golang/internal/logger"
	"react-example/backend-golang/internal/repositories"
	"react-example/backend-golang/internal/usecases"
	"react-example/backend-golang/middleware"
	"react-example/backend-golang/routes"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the IAM Governance API server",
	Run: func(cmd *cobra.Command, args []string) {
		runFiberServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

// @title IAM Governance API
// @version 1.0
// @description This is a sample Identity Governance Suite API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func runFiberServer() {
	// Initialize custom Zap logger
	customLogger.InitLogger()
	if customLogger.Log != nil {
		defer func() {
			_ = customLogger.Log.Sync()
		}()
	}

	customLogger.Log.Info("[IAM-CORE] Bootstrapping Identity Governance Suite via Fiber...")

	db := config.InitDB()
	defer db.Close()

	userRepo := repositories.NewUserRepository(db)
	auditRepo := repositories.NewAuditRepository(db)

	userUsecase := usecases.NewUserUsecase(userRepo, auditRepo)
	auditUsecase := usecases.NewAuditUsecase(auditRepo)
	authUsecase := usecases.NewAuthUsecase(userRepo)
	roleUsecase := usecases.NewRoleUsecase()
	arUsecase := usecases.NewAccessRequestUsecase()
	kycUsecase := usecases.NewKYCUsecase()
	policyUsecase := usecases.NewPolicyUsecase()
	nfUsecase := usecases.NewNotificationUsecase()
	reportUsecase := usecases.NewReportUsecase()

	hc := routes.HandlerContainer{
		UserHandler:         handlers.NewUserHandler(userUsecase),
		AuditHandler:        handlers.NewAuditHandler(auditUsecase),
		AuthHandler:         handlers.NewAuthHandler(authUsecase),
		RoleHandler:         handlers.NewRoleHandler(roleUsecase),
		ARHandler:           handlers.NewAccessRequestHandler(arUsecase),
		KYCHandler:          handlers.NewKYCHandler(kycUsecase),
		PolicyHandler:       handlers.NewPolicyHandler(policyUsecase),
		NotificationHandler: handlers.NewNotificationHandler(nfUsecase),
		ReportHandler:       handlers.NewReportHandler(reportUsecase),
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: false,
		AppName:               "IAM Governance API",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format:     `{"time":"${time}","status":${status},"latency":"${latency}","method":"${method}","path":"${path}","request_id":"${respHeader:X-Request-ID}","error":"${error}"}` + "\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     customLogger.LogWriter,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Content-Type, Authorization",
	}))

	// Register Routes
	routes.RegisterHandlers(app, hc)

	port := config.AppConfig.Port
	
	// Graceful shutdown channel
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			customLogger.Log.Fatal(fmt.Sprintf("[IAM-CORE] CRITICAL: Server crashed: %v", err))
		}
	}()

	<-shutdownSignal
	customLogger.Log.Info("[IAM-CORE] Received shutdown signal, commencing graceful shutdown...")
	if err := app.Shutdown(); err != nil {
		customLogger.Log.Error(fmt.Sprintf("[IAM-CORE] Error during shutdown: %v", err))
	}
}
