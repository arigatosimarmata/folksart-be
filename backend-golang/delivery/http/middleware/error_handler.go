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

// APIError represents the standardized response structure for all error payloads
type APIError struct {
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
	Timestamp time.Time `json:"timestamp"`
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

// HandleError manages unified error classification, log writing, and mapping
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError
	errLabel := "Internal Server Error"
	customMsg := "An unexpected error occurred on the service."

	// 1. Unwrap and inspect the underlying error type
	var customErr *CustomAppError
	if errors.As(err, &customErr) {
		statusCode = customErr.Status
		customMsg = customErr.Message
		errLabel = http.StatusText(statusCode)
		if customErr.Err != nil {
			log.Printf("[IAM-API] Domain error caught: %v | Context: %s", customErr.Err, customErr.Message)
		} else {
			log.Printf("[IAM-API] Generic domain warning caught: %s", customErr.Message)
		}
	} else if errors.Is(err, sql.ErrNoRows) {
		// sql.ErrNoRows -> Resource Not Found
		statusCode = http.StatusNotFound
		errLabel = "Not Found"
		customMsg = "The requested resource could not be found in the governance directory."
		log.Printf("[IAM-API] Database query returned no records: %v", err)
	} else {
		// Identify driver-level and structural database failures dynamically
		errStr := strings.ToLower(err.Error())
		
		if strings.Contains(errStr, "duplicate entry") || strings.Contains(errStr, "error 1062") {
			// MySQL Unique Key Collision
			statusCode = http.StatusConflict
			errLabel = "Conflict"
			customMsg = "A database unique constraint conflict occurred. A record with matching criteria already exists."
			log.Printf("[IAM-API] Unique Constraint Collision: %v", err)
		} else if strings.Contains(errStr, "cannot add or update a child row") || strings.Contains(errStr, "foreign key constraint fails") || strings.Contains(errStr, "error 1452") {
			// MySQL Foreign Key Constraint Collision
			statusCode = http.StatusBadRequest
			errLabel = "Bad Request"
			customMsg = "A relational integrity database constraint failed. Verified matching foreign reference not found."
			log.Printf("[IAM-API] Integrity Constraint Fail: %v", err)
		} else if strings.Contains(errStr, "connection refused") || strings.Contains(errStr, "packets.go") || strings.Contains(errStr, "bad connection") || strings.Contains(errStr, "dial tcp") {
			// Database Host Isolation / Downtime
			statusCode = http.StatusServiceUnavailable
			errLabel = "Service Unavailable"
			customMsg = "The persistent database system is temporarily unavailable. Please retry shortly."
			log.Printf("[IAM-API] RDBMS Host Connection Failure: %v", err)
		} else if strings.Contains(errStr, "syntax error") || strings.Contains(errStr, "sql:") || strings.Contains(errStr, "mysql:") {
			// Structural DB syntax run errors
			statusCode = http.StatusInternalServerError
			errLabel = "Internal Database Error"
			customMsg = "An internal database engine query execution failure occurred."
			log.Printf("[IAM-API] DB Query Syntax/Execution Error: %v", err)
		} else {
			// General system fallback
			log.Printf("[IAM-API] Unhandled server error: %v", err)
			customMsg = err.Error()
		}
	}

	// 2. Generate standardized JSON response
	response := APIError{
		Status:    statusCode,
		Error:     errLabel,
		Message:   customMsg,
		Path:      r.URL.Path,
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
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
