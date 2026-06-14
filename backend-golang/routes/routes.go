package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"react-example/backend-golang/internal/handlers"
	_ "react-example/backend-golang/docs" // Import generated docs
	"time"
)

type HandlerContainer struct {
	UserHandler         *handlers.UserHandler
	AuditHandler        *handlers.AuditHandler
	AuthHandler         *handlers.AuthHandler
	RoleHandler         *handlers.RoleHandler
	ARHandler           *handlers.AccessRequestHandler
	KYCHandler          *handlers.KYCHandler
	PolicyHandler       *handlers.PolicyHandler
	NotificationHandler *handlers.NotificationHandler
	ReportHandler       *handlers.ReportHandler
}

func RegisterHandlers(app *fiber.App, hc HandlerContainer) {
	// Swagger Documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Root API Group
	api := app.Group("/api/v1")

	// Limiters
	userLimiter := limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Hour,
	})

	// 1. Identity Directory Operations
	users := api.Group("/users", userLimiter)
	users.Get("/", hc.UserHandler.ListUsers)
	users.Post("/", hc.UserHandler.EnrollUser)
	users.Patch("/:id", hc.UserHandler.UpdateUser)
	users.Delete("/:id", hc.UserHandler.DeleteUser)

	// KYC
	users.Get("/:id/kyc", hc.KYCHandler.Status)
	users.Post("/:id/kyc", hc.KYCHandler.Submit)
	users.Patch("/:id/kyc", hc.KYCHandler.Review)

	// 2. Authentication & Session
	auth := api.Group("/auth")
	auth.Post("/login", hc.AuthHandler.Login)
	auth.Post("/logout", hc.AuthHandler.Logout)
	auth.Get("/me", hc.AuthHandler.Me)
	
	api.Get("/sessions", hc.AuthHandler.Sessions)

	// 3. Role & Permission
	roles := api.Group("/roles")
	roles.Get("/", hc.RoleHandler.ListRoles)
	roles.Post("/", hc.RoleHandler.CreateRole)
	roles.Post("/assign/:id", hc.RoleHandler.AssignUserRole)

	api.Get("/permissions", hc.RoleHandler.ListPermissions)

	// 4. Access Request & Approval
	ar := api.Group("/access-requests")
	ar.Get("/", hc.ARHandler.List)
	ar.Post("/", hc.ARHandler.Submit)
	ar.Post("/:id/approve", hc.ARHandler.Approve)
	ar.Post("/:id/reject", hc.ARHandler.Reject)

	// 5. Policy Engine
	policies := api.Group("/policies")
	policies.Get("/", hc.PolicyHandler.List)
	policies.Post("/", hc.PolicyHandler.Create)
	policies.Patch("/:id", hc.PolicyHandler.Update)
	policies.Delete("/:id", hc.PolicyHandler.Delete)
	policies.Post("/evaluate", hc.PolicyHandler.Evaluate)

	// 6. Audit Trails
	audit := api.Group("/audit-logs")
	audit.Get("/", hc.AuditHandler.ListAuditLogs)
	audit.Post("/", hc.AuditHandler.CreateLog)
	audit.Post("/sign", hc.AuditHandler.SignLogs)

	// 7. Inter-Service Auth (Internal)
	internal := api.Group("/internal")
	internal.Post("/token", hc.AuthHandler.InternalToken)
	internal.Get("/token/verify", hc.AuthHandler.VerifyInternalToken)

	// 8. Notification & Alerting
	api.Get("/notification-rules", hc.NotificationHandler.ListRules)
	api.Post("/notification-rules", hc.NotificationHandler.CreateRule)
	api.Get("/notifications", hc.NotificationHandler.ListNotifications)

	// 9. Report & Export
	reports := api.Group("/reports")
	reports.Get("/access-summary", hc.ReportHandler.AccessSummary)
	reports.Get("/risk-score-trend", hc.ReportHandler.RiskTrend)

	// 10. API Documentation (Static)
	app.Static("/docs", "./backend-golang/docs")
}
