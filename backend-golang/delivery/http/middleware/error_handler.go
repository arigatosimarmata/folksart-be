package middleware

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// APIError redefined to match PRD: { "success": false, "error": { "code": "...", "message": "..." } }
type APIError struct {
	Success bool       `json:"success"`
	Error   ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// SuccessResponse for standardized success payloads
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

// CustomAppError allows handlers to declare application-specific errors with standard statuses
type CustomAppError struct {
	Status  int
	Message string
	Err     error
}

func (e *CustomAppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%d: %s (%v)", e.Status, e.Message, e.Err)
	}
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

// NewCustomError helper to instantiate application-level errors
func NewCustomError(status int, message string, err error) error {
	return &CustomAppError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

// AppHandler represents standard HTTP handlers with centralized error outputs
type AppHandler func(w http.ResponseWriter, r *http.Request) error

// Adapt turns a central AppHandler into a standard http.HandlerFunc
func Adapt(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			HandleError(w, r, err)
		}
	}
}

// SendJSON helper for standardized success responses
func SendJSON(w http.ResponseWriter, status int, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(SuccessResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

// HandleError manages unified error classification, log writing, and mapping
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError
	errCode := "INTERNAL_SERVER_ERROR"
	customMsg := "An unexpected error occurred on the service."

	// 1. Unwrap and inspect the underlying error type
	var customErr *CustomAppError
	if errors.As(err, &customErr) {
		statusCode = customErr.Status
		customMsg = customErr.Message
		errCode = getErrorCodeFromStatus(statusCode)
		if customErr.Err != nil {
			log.Printf("[IAM-API] Domain error caught: %v | Context: %s", customErr.Err, customErr.Message)
		} else {
			log.Printf("[IAM-API] Generic domain warning caught: %s", customErr.Message)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		statusCode = http.StatusNotFound
		errCode = "NOT_FOUND"
		customMsg = "The requested resource could not be found in the governance directory."
		log.Printf("[IAM-API] Database query returned no records: %v", err)
	} else {
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "duplicate entry") || strings.Contains(errStr, "error 1062") {
			statusCode = http.StatusConflict
			errCode = "CONFLICT"
			customMsg = "A database unique constraint conflict occurred. A record with matching criteria already exists."
		} else if strings.Contains(errStr, "cannot add or update a child row") || strings.Contains(errStr, "foreign key constraint fails") || strings.Contains(errStr, "error 1452") {
			statusCode = http.StatusUnprocessableEntity
			errCode = "VALIDATION_ERROR"
			customMsg = "A relational integrity database constraint failed. Verified matching foreign reference not found."
		} else if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "dial tcp") {
			statusCode = http.StatusServiceUnavailable
			errCode = "SERVICE_UNAVAILABLE"
			customMsg = "The persistent database system is temporarily unavailable. Please retry shortly."
		} else {
			log.Printf("[IAM-API] Unhandled server error: %v", err)
			customMsg = err.Error()
		}
	}

	// 2. Generate standardized JSON response
	response := APIError{
		Success: false,
		Error: ErrorDetail{
			Code:      errCode,
			Message:   customMsg,
			Timestamp: time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}

func getErrorCodeFromStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "BAD_REQUEST"
	case http.StatusUnauthorized:
		return "UNAUTHORIZED"
	case http.StatusForbidden:
		return "FORBIDDEN"
	case http.StatusNotFound:
		return "NOT_FOUND"
	case http.StatusConflict:
		return "CONFLICT"
	case http.StatusUnprocessableEntity:
		return "VALIDATION_ERROR"
	case http.StatusTooManyRequests:
		return "RATE_LIMIT_EXCEEDED"
	default:
		return "INTERNAL_SERVER_ERROR"
	}
}

// RecoveryMiddleware intercepts panics, prevents service crashes, and redirects logs to central handlers
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rcv := recover(); rcv != nil {
				log.Printf("[PANIC-SHIELD] Recovered from deep application panic!\nPanic details: %v\nStack trace:\n%s", rcv, debug.Stack())
				
				var panicErr error
				if e, ok := rcv.(error); ok {
					panicErr = e
				} else {
					panicErr = fmt.Errorf("%v", rcv)
				}
				
				HandleError(w, r, NewCustomError(http.StatusInternalServerError, "An unexpected critical runtime panic occurred on the server.", panicErr))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
