package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"react-example/backend-golang/internal/ctxutil"
)

// RequestID returns a middleware that injects a request ID into the request context and headers.
func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Read from header or generate a new UUID
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		// Echo the request ID back in the response headers
		c.Set("X-Request-ID", reqID)

		// Set the request ID in the Go standard Context (UserContext)
		ctx := ctxutil.WithRequestID(c.UserContext(), reqID)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
