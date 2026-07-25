package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"folksart-be/backend-golang/internal/handlers"
	"folksart-be/backend-golang/middleware"
	_ "folksart-be/backend-golang/docs" // Import generated docs
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

	// Public Authentication Endpoint
	auth := api.Group("/auth")
	auth.Post("/login", hc.AuthHandler.Login)
	auth.Post("/register", hc.AuthHandler.RegisterCustomer)

	// Protected Route Group (Requires Bearer JWT Token)
	protected := api.Group("", middleware.JWTAuth())

	// Protected Auth Endpoints & Sessions
	protectedAuth := protected.Group("/auth")
	protectedAuth.Post("/logout", hc.AuthHandler.Logout)
	protectedAuth.Get("/me", hc.AuthHandler.Me)
	protected.Get("/sessions", hc.AuthHandler.Sessions)

	// 1. Identity Directory Operations
	users := protected.Group("/users", userLimiter)
	users.Get("/", hc.UserHandler.ListUsers)
	users.Post("/", middleware.RBAC("Administrator"), hc.UserHandler.EnrollUser)
	users.Patch("/:id", hc.UserHandler.UpdateUser)
	users.Delete("/:id", middleware.RBAC("Administrator"), hc.UserHandler.DeleteUser)

	// KYC
	users.Get("/:id/kyc", hc.KYCHandler.Status)
	users.Post("/:id/kyc", hc.KYCHandler.Submit)
	users.Patch("/:id/kyc", middleware.RBAC("Administrator"), hc.KYCHandler.Review)

	// 2. Role & Permission
	roles := protected.Group("/roles")
	roles.Get("/", hc.RoleHandler.ListRoles)
	roles.Post("/", middleware.RBAC("Administrator"), hc.RoleHandler.CreateRole)
	roles.Post("/assign/:id", middleware.RBAC("Administrator"), hc.RoleHandler.AssignUserRole)
	protected.Get("/permissions", hc.RoleHandler.ListPermissions)

	// 3. Access Request & Approval
	ar := protected.Group("/access-requests")
	ar.Get("/", hc.ARHandler.List)
	ar.Post("/", hc.ARHandler.Submit)
	ar.Post("/:id/approve", middleware.RBAC("Administrator"), hc.ARHandler.Approve)
	ar.Post("/:id/reject", middleware.RBAC("Administrator"), hc.ARHandler.Reject)

	// 4. Policy Engine
	policies := protected.Group("/policies")
	policies.Get("/", hc.PolicyHandler.List)
	policies.Post("/", middleware.RBAC("Administrator"), hc.PolicyHandler.Create)
	policies.Patch("/:id", middleware.RBAC("Administrator"), hc.PolicyHandler.Update)
	policies.Delete("/:id", middleware.RBAC("Administrator"), hc.PolicyHandler.Delete)
	policies.Post("/evaluate", hc.PolicyHandler.Evaluate)

	// 5. Audit Trails
	audit := protected.Group("/audit-logs")
	audit.Get("/", hc.AuditHandler.ListAuditLogs)
	audit.Post("/", hc.AuditHandler.CreateLog)
	audit.Post("/sign", middleware.RBAC("Administrator"), hc.AuditHandler.SignLogs)

	// 6. Inter-Service Auth (Internal)
	internal := protected.Group("/internal")
	internal.Post("/token", hc.AuthHandler.InternalToken)
	internal.Get("/token/verify", hc.AuthHandler.VerifyInternalToken)

	// 7. Notification & Alerting
	protected.Get("/notification-rules", hc.NotificationHandler.ListRules)
	protected.Post("/notification-rules", middleware.RBAC("Administrator"), hc.NotificationHandler.CreateRule)
	protected.Get("/notifications", hc.NotificationHandler.ListNotifications)

	// 8. Report & Export
	reports := protected.Group("/reports")
	reports.Get("/access-summary", hc.ReportHandler.AccessSummary)
	reports.Get("/risk-score-trend", hc.ReportHandler.RiskTrend)

	// 9. API Documentation (Static)
	app.Static("/docs", "./backend-golang/docs")
}
