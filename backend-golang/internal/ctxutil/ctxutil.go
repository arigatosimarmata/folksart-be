package ctxutil

import (
	"context"
)

type ctxKey string

const reqIDKey ctxKey = "request_id"

// WithRequestID puts a request ID into the given context.
func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, reqIDKey, reqID)
}

// GetRequestID gets a request ID from the given context, or empty string if not found.
func GetRequestID(ctx context.Context) string {
	val := ctx.Value(reqIDKey)
	if val == nil {
		return ""
	}
	str, ok := val.(string)
	if !ok {
		return ""
	}
	return str
}
