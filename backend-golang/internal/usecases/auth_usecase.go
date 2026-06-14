package usecases

import (
	"context"
	"fmt"
	"time"

	"react-example/backend-golang/internal/domain"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(ur domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{userRepo: ur}
}

func (u *authUsecase) Login(ctx context.Context, email, password string) (*domain.TokenResponse, error) {
	// Mock login for boilerplate
	if email == "" || password == "" {
		return nil, fmt.Errorf("credentials required")
	}

	return &domain.TokenResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		TokenType:   "Bearer",
		ExpiresIn:   3600,
		User: &domain.IAMUser{
			ID:    "usr-admin",
			Name:  "Admin user",
			Email: email,
			Role:  "Administrator",
		},
	}, nil
}

func (u *authUsecase) Logout(ctx context.Context, token string) error {
	return nil
}

func (u *authUsecase) Refresh(ctx context.Context, token string) (*domain.TokenResponse, error) {
	return nil, nil
}

func (u *authUsecase) Me(ctx context.Context, token string) (*domain.IAMUser, error) {
	return &domain.IAMUser{
		ID:    "usr-admin",
		Name:  "Admin user",
		Email: "admin@example.com",
		Role:  "Administrator",
	}, nil
}

func (u *authUsecase) GetSessions(ctx context.Context, userID string) ([]domain.Session, error) {
	return []domain.Session{
		{
			ID:           "sess-1",
			Device:       "MacBook Pro - Chrome",
			IPAddress:    "192.168.1.1",
			Location:     "Jakarta, Indonesia",
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			LastActiveAt: time.Now(),
			IsCurrent:    true,
		},
	}, nil
}

func (u *authUsecase) TerminateSession(ctx context.Context, userID, sessionID string) error {
	return nil
}
