package usecase

import (
	"context"
	"errors"
	"time"

	"react-example/backend-golang/domain"
)

type authUsecase struct {
	userRepo domain.UserRepository
}

func NewAuthUsecase(ur domain.UserRepository) domain.AuthUsecase {
	return &authUsecase{userRepo: ur}
}

func (a *authUsecase) Login(ctx context.Context, email, password string) (*domain.TokenResponse, error) {
	// Mock: verify password is 'password'
	if password != "password" {
		return nil, errors.New("invalid credentials")
	}

	return &domain.TokenResponse{
		AccessToken: "tok_" + time.Now().Format("20060102150405"),
		TokenType:   "Bearer",
		ExpiresIn:   86400,
		User: &domain.IAMUser{
			ID:    "usr_01HXYZ123456",
			Name:  "John Doe",
			Email: email,
			Role:  "admin",
		},
	}, nil
}

func (a *authUsecase) Logout(ctx context.Context, token string) error {
	return nil
}

func (a *authUsecase) Refresh(ctx context.Context, token string) (*domain.TokenResponse, error) {
	return &domain.TokenResponse{
		AccessToken: "tok_ref_" + time.Now().Format("20060102150405"),
		TokenType:   "Bearer",
		ExpiresIn:   86400,
	}, nil
}

func (a *authUsecase) Me(ctx context.Context, token string) (*domain.IAMUser, error) {
	return &domain.IAMUser{
		ID:         "usr_01HXYZ123456",
		Email:      "john.doe@company.com",
		Name:       "John Doe",
		Department: "Engineering",
		Role:       "admin",
		Status:     "active",
		RiskScore:  12,
		KYCStatus:  "verified",
		CreatedAt:  time.Now().AddDate(-1, 0, 0),
	}, nil
}

func (a *authUsecase) GetSessions(ctx context.Context, userID string) ([]domain.Session, error) {
	return []domain.Session{
		{
			ID:           "sess_01ABC",
			Device:       "Chrome on macOS",
			IPAddress:    "192.168.1.100",
			Location:     "Bandung, ID",
			CreatedAt:    time.Now().Add(-24 * time.Hour),
			LastActiveAt: time.Now(),
			IsCurrent:    true,
		},
	}, nil
}

func (a *authUsecase) TerminateSession(ctx context.Context, userID, sessionID string) error {
	return nil
}
