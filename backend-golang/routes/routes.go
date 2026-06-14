package routes

import (
	"net/http"
	"strings"
	"time"

	"react-example/backend-golang/internal/handlers"
	"react-example/backend-golang/middleware"
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

func RegisterHandlers(hc HandlerContainer) {
	userLimiter := middleware.NewRateLimiter(5.0, 10.0, 1*time.Hour)
	auditLimiter := middleware.NewRateLimiter(3.0, 5.0, 30*time.Minute)
	authLimiter := middleware.NewRateLimiter(10.0, 20.0, 1*time.Hour)

	// 1. Identity Directory Operations
	http.HandleFunc("/api/v1/users", middleware.LimitMiddleware(userLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.UserHandler.ListUsers(w, r)
		} else if r.Method == http.MethodPost {
			return hc.UserHandler.EnrollUser(w, r)
		}
		return nil
	})))

	http.HandleFunc("/api/v1/users/", middleware.LimitMiddleware(userLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if strings.Contains(r.URL.Path, "/kyc") {
			if r.Method == http.MethodGet {
				return hc.KYCHandler.Status(w, r)
			} else if r.Method == http.MethodPost {
				return hc.KYCHandler.Submit(w, r)
			} else if r.Method == http.MethodPatch {
				return hc.KYCHandler.Review(w, r)
			}
		}
		if r.Method == http.MethodPatch {
			return hc.UserHandler.UpdateUser(w, r)
		} else if r.Method == http.MethodDelete {
			return hc.UserHandler.DeleteUser(w, r)
		}
		return nil
	})))

	// 2. Authentication & Session
	http.HandleFunc("/api/v1/auth/login", middleware.LimitMiddleware(authLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodPost {
			return hc.AuthHandler.Login(w, r)
		}
		return nil
	})))

	http.HandleFunc("/api/v1/auth/logout", middleware.Adapt(hc.AuthHandler.Logout))
	http.HandleFunc("/api/v1/auth/me", middleware.Adapt(hc.AuthHandler.Me))
	http.HandleFunc("/api/v1/sessions", middleware.Adapt(hc.AuthHandler.Sessions))

	// 3. Role & Permission
	http.HandleFunc("/api/v1/roles", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.RoleHandler.ListRoles(w, r)
		} else if r.Method == http.MethodPost {
			return hc.RoleHandler.CreateRole(w, r)
		}
		return nil
	}))

	http.HandleFunc("/api/v1/permissions", middleware.Adapt(hc.RoleHandler.ListPermissions))

	// 4. Access Request & Approval
	http.HandleFunc("/api/v1/access-requests", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.ARHandler.List(w, r)
		} else if r.Method == http.MethodPost {
			return hc.ARHandler.Submit(w, r)
		}
		return nil
	}))

	http.HandleFunc("/api/v1/access-requests/", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if strings.HasSuffix(r.URL.Path, "/approve") {
			return hc.ARHandler.Approve(w, r)
		} else if strings.HasSuffix(r.URL.Path, "/reject") {
			return hc.ARHandler.Reject(w, r)
		}
		return nil
	}))

	// 5. Policy Engine
	http.HandleFunc("/api/v1/policies", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.PolicyHandler.List(w, r)
		} else if r.Method == http.MethodPost {
			return hc.PolicyHandler.Create(w, r)
		}
		return nil
	}))

	http.HandleFunc("/api/v1/policies/evaluate", middleware.Adapt(hc.PolicyHandler.Evaluate))

	// 6. Audit Trails
	http.HandleFunc("/api/v1/audit-logs", middleware.LimitMiddleware(auditLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.AuditHandler.ListAuditLogs(w, r)
		} else if r.Method == http.MethodPost {
			return hc.AuditHandler.CreateLog(w, r)
		}
		return nil
	})))

	http.HandleFunc("/api/v1/audit-logs/sign", middleware.Adapt(hc.AuditHandler.SignLogs))

	// 7. Inter-Service Auth (Internal)
	http.HandleFunc("/api/v1/internal/token", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodPost {
			return hc.AuthHandler.InternalToken(w, r)
		}
		return nil
	}))

	http.HandleFunc("/api/v1/internal/token/verify", middleware.Adapt(hc.AuthHandler.VerifyInternalToken))

	// 8. Notification & Alerting
	http.HandleFunc("/api/v1/notification-rules", middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodGet {
			return hc.NotificationHandler.ListRules(w, r)
		} else if r.Method == http.MethodPost {
			return hc.NotificationHandler.CreateRule(w, r)
		}
		return nil
	}))

	http.HandleFunc("/api/v1/notifications", middleware.Adapt(hc.NotificationHandler.ListNotifications))

	// 8. Report & Export
	http.HandleFunc("/api/v1/reports/access-summary", middleware.Adapt(hc.ReportHandler.AccessSummary))
	http.HandleFunc("/api/v1/reports/risk-score-trend", middleware.Adapt(hc.ReportHandler.RiskTrend))

	// 9. API Documentation
	fs := http.FileServer(http.Dir("./backend-golang/docs"))
	http.Handle("/docs/", http.StripPrefix("/docs/", fs))
}
