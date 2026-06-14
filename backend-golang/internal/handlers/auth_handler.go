package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"react-example/backend-golang/httputil"
	"react-example/backend-golang/internal/domain"
	"react-example/backend-golang/internal/dto"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(au domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: au}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	resp, err := h.authUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	dtoResp := dto.TokenResponse{
		AccessToken: resp.AccessToken,
		TokenType:   resp.TokenType,
		ExpiresIn:   resp.ExpiresIn,
	}

	if resp.User != nil {
		dtoResp.User = &dto.UserResponse{
			ID:         resp.User.ID,
			Name:       resp.User.Name,
			Username:   resp.User.Username,
			Email:      resp.User.Email,
			Role:       resp.User.Role,
			Status:     resp.User.Status,
			KYCStatus:  resp.User.KYCStatus,
			Department: resp.User.Department,
			RiskScore:  resp.User.RiskScore,
			CreatedAt:  resp.User.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	httputil.WriteSuccessResponse(w, "Login successful", dtoResp, nil)
	return nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := h.authUsecase.Logout(r.Context(), token)
	if err != nil {
		return err
	}

	httputil.WriteSuccessResponse(w, "Logout successful", nil, nil)
	return nil
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	user, err := h.authUsecase.Me(r.Context(), token)
	if err != nil {
		return err
	}

	res := dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Username:   user.Username,
		Email:      user.Email,
		Role:       user.Role,
		Status:     user.Status,
		KYCStatus:  user.KYCStatus,
		Department: user.Department,
		RiskScore:  user.RiskScore,
		CreatedAt:  user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	httputil.WriteSuccessResponse(w, "Success", res, nil)
	return nil
}

func (h *AuthHandler) Sessions(w http.ResponseWriter, r *http.Request) error {
	sessions, err := h.authUsecase.GetSessions(r.Context(), "usr_current")
	if err != nil {
		return err
	}

	dtoSessions := make([]dto.SessionResponse, 0)
	for _, s := range sessions {
		dtoSessions = append(dtoSessions, dto.SessionResponse{
			ID:           s.ID,
			Device:       s.Device,
			IPAddress:    s.IPAddress,
			Location:     s.Location,
			CreatedAt:    s.CreatedAt.Format("2006-01-02 15:04:05"),
			LastActiveAt: s.LastActiveAt.Format("2006-01-02 15:04:05"),
			IsCurrent:    s.IsCurrent,
		})
	}

	httputil.WriteSuccessResponse(w, "Success", dtoSessions, nil)
	return nil
}

func (h *AuthHandler) InternalToken(w http.ResponseWriter, r *http.Request) error {
	httputil.WriteSuccessResponse(w, "Internal token generated", map[string]interface{}{
		"access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"token_type":   "Bearer",
		"expires_in":   3600,
		"issued_to":    "reporting-service",
	}, nil)
	return nil
}

func (h *AuthHandler) VerifyInternalToken(w http.ResponseWriter, r *http.Request) error {
	httputil.WriteSuccessResponse(w, "Token verified", map[string]interface{}{
		"valid":     true,
		"issued_to": "reporting-service",
	}, nil)
	return nil
}
