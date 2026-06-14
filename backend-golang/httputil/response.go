package httputil

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"react-example/backend-golang/errs"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

func WriteSuccessResponse(c *fiber.Ctx, message string, data interface{}, meta interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Code:    "00",
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func WriteErrorResponse(c *fiber.Ctx, err error) error {
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

	return c.Status(statusCode).JSON(Response{
		Code:    code,
		Message: message,
	})
}

func WriteValidationErrorResponse(c *fiber.Ctx, validationErrors interface{}) error {
	return c.Status(http.StatusBadRequest).JSON(Response{
		Code:    "40",
		Message: "Validation Failed",
		Data:    validationErrors,
	})
}
