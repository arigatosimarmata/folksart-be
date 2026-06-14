package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"react-example/backend-golang/config"
	"react-example/backend-golang/internal/handlers"
	"react-example/backend-golang/internal/repositories"
	"react-example/backend-golang/internal/usecases"
	"react-example/backend-golang/middleware"
	"react-example/backend-golang/routes"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the IAM Governance API server",
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

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

func runServer() {
	log.Println("[IAM-CORE] Bootstrapping Identity Governance Suite...")

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

	routes.RegisterHandlers(hc)

	port := config.AppConfig.Port
	srv := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: corsMiddleware(middleware.RecoveryMiddleware(http.DefaultServeMux)),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("[IAM-CORE] Service listening dynamically at http://0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Fatalf("[IAM-CORE] CRITICAL: Server crashed: %v", err)
	case sig := <-shutdownSignal:
		log.Printf("[IAM-CORE] Received signal %v, commencing graceful shutdown...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("[IAM-CORE] Graceful server drainage failed: %v", err)
			_ = srv.Close()
		}
	}
}
