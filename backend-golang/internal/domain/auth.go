package domain

import (
	"context"
	"time"
)

type TokenResponse struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	User        *IAMUser
}

type Session struct {
	ID           string
	Device       string
	IPAddress    string
	Location     string
	CreatedAt    time.Time
	LastActiveAt time.Time
	IsCurrent    bool
}

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*TokenResponse, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, token string) (*TokenResponse, error)
	Me(ctx context.Context, token string) (*IAMUser, error)
	GetSessions(ctx context.Context, userID string) ([]Session, error)
	TerminateSession(ctx context.Context, userID, sessionID string) error
}
