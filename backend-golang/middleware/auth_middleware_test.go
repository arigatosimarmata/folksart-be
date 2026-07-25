package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"folksart-be/backend-golang/config"
	"folksart-be/backend-golang/internal/security"
)

func TestJWTAuthAndRBAC(t *testing.T) {
	config.AppConfig.JWTSecretKey = "test-secret"
	validAdminToken, _ := security.GenerateToken("usr-1", "Administrator")
	validUserToken, _ := security.GenerateToken("usr-2", "Viewer")

	app := fiber.New()
	app.Use("/protected", JWTAuth())
	app.Get("/protected/data", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})
	app.Delete("/protected/admin", RBAC("Administrator"), func(c *fiber.Ctx) error {
		return c.SendString("admin ok")
	})

	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected/data", nil)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected/data", nil)
		req.Header.Set("Authorization", "Token some-invalid-format")
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", resp.StatusCode)
		}
	})

	t.Run("Valid Token Access to Data", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/protected/data", nil)
		req.Header.Set("Authorization", "Bearer "+validUserToken)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", resp.StatusCode)
		}
	})

	t.Run("Forbidden RBAC Access", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/protected/admin", nil)
		req.Header.Set("Authorization", "Bearer "+validUserToken)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("expected 403 Forbidden for Viewer role, got %d", resp.StatusCode)
		}
	})

	t.Run("Allowed RBAC Access for Administrator", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/protected/admin", nil)
		req.Header.Set("Authorization", "Bearer "+validAdminToken)
		resp, _ := app.Test(req)
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 OK for Administrator role, got %d", resp.StatusCode)
		}
	})
}
