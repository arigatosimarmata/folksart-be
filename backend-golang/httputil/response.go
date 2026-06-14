package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"react-example/backend-golang/errs"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func WriteSuccessResponse(w http.ResponseWriter, message string, data interface{}, meta interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Code:    "00",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	code := "99"
	message := err.Error()

	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		code = "44"
	case errors.Is(err, errs.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		code = "41"
	case errors.Is(err, errs.ErrForbidden):
		statusCode = http.StatusForbidden
		code = "43"
	case errors.Is(err, errs.ErrBadRequest) || errors.Is(err, errs.ErrValidation):
		statusCode = http.StatusBadRequest
		code = "40"
	case errors.Is(err, errs.ErrConflict):
		statusCode = http.StatusConflict
		code = "49"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Code:    code,
		Message: message,
	})
}

func WriteValidationErrorResponse(w http.ResponseWriter, validationErrors interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{
		Code:    "40",
		Message: "Validation Failed",
		Data:    validationErrors,
	})
}
