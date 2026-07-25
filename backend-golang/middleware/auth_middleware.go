package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"folksart-be/backend-golang/httputil"
	"folksart-be/backend-golang/internal/security"
)

// JWTAuth middleware verifies the Bearer JWT token in the Authorization header
func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(httputil.Response{
				Code:    "41",
				Message: "Unauthorized: missing Authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(http.StatusUnauthorized).JSON(httputil.Response{
				Code:    "41",
				Message: "Unauthorized: invalid Authorization format, expected 'Bearer <token>'",
			})
		}

		claims, err := security.ValidateToken(parts[1])
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(httputil.Response{
				Code:    "41",
				Message: "Unauthorized: " + err.Error(),
			})
		}

		// Save claims in Fiber locals for downstream handlers
		c.Locals("userId", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

// RBAC middleware enforces Role-Based Access Control for specified roles
func RBAC(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(http.StatusForbidden).JSON(httputil.Response{
				Code:    "43",
				Message: "Forbidden: role not defined in request context",
			})
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(http.StatusForbidden).JSON(httputil.Response{
			Code:    "43",
			Message: "Forbidden: insufficient permissions for role '" + userRole + "'",
		})
	}
}
