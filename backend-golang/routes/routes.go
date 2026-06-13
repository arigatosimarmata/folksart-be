package routes

import (
	"net/http"
	"time"

	delivery "react-example/backend-golang/delivery/http"
	"react-example/backend-golang/delivery/http/middleware"
)

// RegisterHandlers maps corporate endpoints under the /api/v1 namespace using Clean Architecture controllers
func RegisterHandlers(userHandler *delivery.UserHandler, auditHandler *delivery.AuditHandler) {
	// Initialize token bucket rate limiters for separate resources
	userLimiter := middleware.NewRateLimiter(5.0, 10.0, 1*time.Hour)       // 5 reqs/sec, burst capacity of 10
	auditLimiter := middleware.NewRateLimiter(3.0, 5.0, 30*time.Minute)   // 3 reqs/sec, burst capacity of 5

	// 1. Identity Directory Operations
	http.HandleFunc("/api/v1/users", middleware.LimitMiddleware(userLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodOptions {
			return nil // Handled at corsMiddleware level
		}
		if r.Method == http.MethodGet {
			return userHandler.ListUsers(w, r)
		} else if r.Method == http.MethodPost {
			return userHandler.EnrollUser(w, r)
		} else {
			return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		}
	})))

	// Handle ID sub-routing: /api/v1/users/usr-XXXX
	http.HandleFunc("/api/v1/users/", middleware.LimitMiddleware(userLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodOptions {
			return nil // Handled at corsMiddleware level
		}
		if r.Method == http.MethodPatch {
			return userHandler.UpdateUser(w, r)
		} else if r.Method == http.MethodDelete {
			return userHandler.DeleteUser(w, r)
		} else {
			return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		}
	})))

	// 2. Audit Trails Operations
	http.HandleFunc("/api/v1/audit-logs", middleware.LimitMiddleware(auditLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodOptions {
			return nil // Handled at corsMiddleware level
		}
		if r.Method == http.MethodGet {
			return auditHandler.ListAuditLogs(w, r)
		} else if r.Method == http.MethodPost {
			return auditHandler.CreateLog(w, r)
		} else {
			return middleware.NewCustomError(http.StatusMethodNotAllowed, "Method Not Allowed", nil)
		}
	})))

	// 3. CSV Dataset Export Engine
	http.HandleFunc("/api/v1/export/csv", middleware.LimitMiddleware(userLimiter, middleware.Adapt(func(w http.ResponseWriter, r *http.Request) error {
		if r.Method == http.MethodOptions {
			return nil // Handled at corsMiddleware level
		}
		return userHandler.ExportCSV(w, r)
	})))
}
