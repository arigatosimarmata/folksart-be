package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/delivery/http/middleware"
	"react-example/backend-golang/domain"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(au domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: au}
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var payload LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return middleware.NewCustomError(http.StatusBadRequest, "Invalid payload", err)
	}

	resp, err := h.authUsecase.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		return middleware.NewCustomError(http.StatusUnauthorized, err.Error(), nil)
	}

	middleware.SendJSON(w, http.StatusOK, resp, nil)
	return nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := h.authUsecase.Logout(r.Context(), token)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, map[string]string{"message": "Logout successful"}, nil)
	return nil
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := h.authUsecase.Me(r.Context(), token)
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, user, nil)
	return nil
}

func (h *AuthHandler) Sessions(w http.ResponseWriter, r *http.Request) error {
	// In real app, get userID from context/token
	sessions, err := h.authUsecase.GetSessions(r.Context(), "usr_current")
	if err != nil {
		return err
	}

	middleware.SendJSON(w, http.StatusOK, sessions, nil)
	return nil
}

func (h *AuthHandler) InternalToken(w http.ResponseWriter, r *http.Request) error {
	middleware.SendJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"issued_to":    "reporting-service",
	}, nil)
	return nil
}

func (h *AuthHandler) VerifyInternalToken(w http.ResponseWriter, r *http.Request) error {
	middleware.SendJSON(w, http.StatusOK, map[string]interface{}{
		"valid":     true,
		"issued_to": "reporting-service",
	}, nil)
	return nil
}
