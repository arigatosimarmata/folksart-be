package ctxutil

import "context"

type contextKey string

const RequestIDKey contextKey = "request_id"

// WithRequestID returns a new context with the request ID attached.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID returns the request ID from the context if it exists.
func GetRequestID(ctx context.Context) string {
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}
