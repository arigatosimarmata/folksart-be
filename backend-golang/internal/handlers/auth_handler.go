package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"folksart-be/backend-golang/httputil"
	"folksart-be/backend-golang/internal/domain"
	"folksart-be/backend-golang/internal/dto"
)

type AuthHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(au domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: au}
}

// Login godoc
// @Summary Authenticate principal
// @Description Log in with email and password to receive a JWT access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} httputil.Response{data=dto.TokenResponse}
// @Failure 401 {object} httputil.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	resp, err := h.authUsecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	dtoResp := dto.TokenResponse{
		AccessToken: resp.AccessToken,
		TokenType:   resp.TokenType,
		ExpiresIn:   resp.ExpiresIn,
	}

	if resp.User != nil {
		dtoResp.User = toUserDTO(resp.User)
	}

	return httputil.WriteSuccessResponse(c, "Login successful", dtoResp, nil)
}

// RegisterCustomer godoc
// @Summary Customer CIAM self-service enrollment
// @Description Register a new external customer identity and get token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterCustomerRequest true "Customer Registration Data"
// @Success 201 {object} httputil.Response{data=dto.TokenResponse}
// @Failure 400 {object} httputil.Response
// @Router /auth/register [post]
func (h *AuthHandler) RegisterCustomer(c *fiber.Ctx) error {
	var req dto.RegisterCustomerRequest
	if err := c.BodyParser(&req); err != nil {
		return httputil.WriteErrorResponse(c, err)
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return httputil.WriteErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Name, email, and password are required"))
	}
	resp, err := h.authUsecase.RegisterCustomer(c.Context(), req.Name, req.Email, req.Password, req.Phone)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}
	dtoResp := dto.TokenResponse{
		AccessToken: resp.AccessToken,
		TokenType:   resp.TokenType,
		ExpiresIn:   resp.ExpiresIn,
	}
	if resp.User != nil {
		dtoResp.User = toUserDTO(resp.User)
	}
	return httputil.WriteSuccessResponse(c, "Customer registration successful", dtoResp, nil)
}

// Logout godoc
// @Summary Terminate session
// @Description Invalidate the current access token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} httputil.Response
// @Security ApiKeyAuth
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	err := h.authUsecase.Logout(c.Context(), token)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Logout successful", nil, nil)
}

// Me godoc
// @Summary Get current principal info
// @Description Retrieve details of the currently authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} httputil.Response{data=dto.UserResponse}
// @Security ApiKeyAuth
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	user, err := h.authUsecase.Me(c.Context(), token)
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
	}

	return httputil.WriteSuccessResponse(c, "Success", toUserDTO(user), nil)
}

func (h *AuthHandler) Sessions(c *fiber.Ctx) error {
	sessions, err := h.authUsecase.GetSessions(c.Context(), "usr_current")
	if err != nil {
		return httputil.WriteErrorResponse(c, err)
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

	return httputil.WriteSuccessResponse(c, "Success", dtoSessions, nil)
}

func (h *AuthHandler) InternalToken(c *fiber.Ctx) error {
	return httputil.WriteSuccessResponse(c, "Internal token generated", map[string]interface{}{
		"accessToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
		"tokenType":   "Bearer",
		"expiresIn":   3600,
		"issuedTo":    "reporting-service",
	}, nil)
}

func (h *AuthHandler) VerifyInternalToken(c *fiber.Ctx) error {
	return httputil.WriteSuccessResponse(c, "Token verified", map[string]interface{}{
		"valid":    true,
		"issuedTo": "reporting-service",
	}, nil)
}
