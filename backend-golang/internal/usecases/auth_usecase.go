package usecases

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"folksart-be/backend-golang/internal/domain"
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
			ID:         "usr-admin",
			Name:       "Admin user",
			Email:      email,
			Role:       "Administrator",
			UserType:   "workforce",
			TenantID:   "internal-org",
			IsVerified: true,
		},
	}, nil
}

func (u *authUsecase) RegisterCustomer(ctx context.Context, name, email, password, phone string) (*domain.TokenResponse, error) {
	if name == "" || email == "" || password == "" {
		return nil, fmt.Errorf("insufficient registration credentials provided")
	}

	cleanUsername := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	if len(cleanUsername) > 15 {
		cleanUsername = cleanUsername[:15]
	}

	rand.Seed(time.Now().UnixNano())
	randomID := fmt.Sprintf("cust-%d", 10000+rand.Intn(90000))
	createdAt := time.Now()

	newUser := &domain.IAMUser{
		ID:         randomID,
		Name:       name,
		Username:   cleanUsername,
		Email:      email,
		Phone:      phone,
		Role:       "VerifiedCustomer",
		Status:     "Active",
		KYCStatus:  "Pending",
		Department: "External",
		RiskScore:  5,
		MFAEnabled: false,
		UserType:   "customer",
		TenantID:   "public-ciam",
		IsVerified: false,
		CreatedAt:  createdAt,
	}

	if u.userRepo != nil {
		_ = u.userRepo.Store(ctx, newUser)
	}

	return &domain.TokenResponse{
		AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkRFTU8tQ1VTVE9NRVJ...",
		TokenType:   "Bearer",
		ExpiresIn:   86400,
		User:        newUser,
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
		ID:         "usr-admin",
		Name:       "Admin user",
		Email:      "admin@example.com",
		Role:       "Administrator",
		UserType:   "workforce",
		TenantID:   "internal-org",
		IsVerified: true,
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
