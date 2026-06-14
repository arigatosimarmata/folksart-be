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
	delivery "react-example/backend-golang/delivery/http"
	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/repository/mysql"
	"react-example/backend-golang/routes"
	"react-example/backend-golang/usecase"
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

	userRepo := mysql.NewMySQLUserRepository(db)
	auditRepo := mysql.NewMySQLAuditRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo, auditRepo)
	auditUsecase := usecase.NewAuditUsecase(auditRepo)
	authUsecase := usecase.NewAuthUsecase(userRepo)
	roleUsecase := usecase.NewRoleUsecase()
	arUsecase := usecase.NewAccessRequestUsecase()
	kycUsecase := usecase.NewKYCUsecase()
	policyUsecase := usecase.NewPolicyUsecase()
	nfUsecase := usecase.NewNotificationUsecase()
	reportUsecase := usecase.NewReportUsecase()

	hc := routes.HandlerContainer{
		UserHandler:         delivery.NewUserHandler(userUsecase),
		AuditHandler:        delivery.NewAuditHandler(auditUsecase),
		AuthHandler:         delivery.NewAuthHandler(authUsecase),
		RoleHandler:         delivery.NewRoleHandler(roleUsecase),
		ARHandler:           delivery.NewAccessRequestHandler(arUsecase),
		KYCHandler:          delivery.NewKYCHandler(kycUsecase),
		PolicyHandler:       delivery.NewPolicyHandler(policyUsecase),
		NotificationHandler: delivery.NewNotificationHandler(nfUsecase),
		ReportHandler:       delivery.NewReportHandler(reportUsecase),
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
