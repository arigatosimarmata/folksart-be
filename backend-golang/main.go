package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"react-example/backend-golang/config"
	delivery "react-example/backend-golang/delivery/http"
	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/repository/mysql"
	"react-example/backend-golang/routes"
	"react-example/backend-golang/usecase"
	"go.uber.org/zap"
)

// corsMiddleware establishes access controls so React frontend browsers can communicate seamlessly
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize structured logger with 5 days rotation
	logger := config.InitLogger(config.LoggerConfig{
		LogPath:    "app.log",
		MaxAgeDays: 5,
	})
	defer logger.Sync()

	logger.Info("[IAM-CORE] Bootstrapping Identity Governance Suite via Clean Architecture...")

	// 1. Establish MySQL RDBMS Connection pool
	db := config.InitDB()
	defer func() {
		logger.Info("[IAM-CORE] Terminating MySQL database pool connection...")
		db.Close()
	}()

	// 2. Setup repositories (Infrastructure Layer)
	userRepo := mysql.NewMySQLUserRepository(db)
	auditRepo := mysql.NewMySQLAuditRepository(db)

	// 3. Setup core interactors (Usecase/Business Layer)
	userUsecase := usecase.NewUserUsecase(userRepo, auditRepo)
	auditUsecase := usecase.NewAuditUsecase(auditRepo)

	// 4. Setup delivery controllers (HTTP Handler Layer)
	userHandler := delivery.NewUserHandler(userUsecase)
	auditHandler := delivery.NewAuditHandler(auditUsecase)

	// 5. Wire routing patterns
	routes.RegisterHandlers(userHandler, auditHandler)

	// 6. Discover target bind port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Baseline standalone port for Go service
	}

	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: corsMiddleware(middleware.RecoveryMiddleware(http.DefaultServeMux)),
	}

	// Channel to capture listen errors
	serverErrors := make(chan error, 1)

	// Direct server listener instantiation in a background thread
	go func() {
		logger.Info("[IAM-CORE] Service listening dynamically", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Signal interception channel definition
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	// Block main thread until signal notification or critical bind error returns
	select {
	case err := <-serverErrors:
		logger.Fatal("[IAM-CORE] CRITICAL: Server crashed during runtime bind", zap.Error(err))

	case sig := <-shutdownSignal:
		logger.Info("[IAM-CORE] Received signal, commencing graceful shutdown routine...", zap.String("signal", sig.String()))

		// Enforce a maximum context timeout threshold of 15 seconds for outstanding client request drainage
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("[IAM-CORE] Graceful server drainage failed, forcing termination", zap.Error(err))
			_ = srv.Close()
		}
	}
}
