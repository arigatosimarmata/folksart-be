package domain

import (
	"context"
	"time"
)

type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	User        *IAMUser  `json:"user,omitempty"`
}

type Session struct {
	ID           string    `json:"id"`
	Device       string    `json:"device"`
	IPAddress    string    `json:"ip_address"`
	Location     string    `json:"location"`
	CreatedAt    time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at"`
	IsCurrent    bool      `json:"is_current"`
}

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*TokenResponse, error)
	Logout(ctx context.Context, token string) error
	Refresh(ctx context.Context, token string) (*TokenResponse, error)
	Me(ctx context.Context, token string) (*IAMUser, error)
	GetSessions(ctx context.Context, userID string) ([]Session, error)
	TerminateSession(ctx context.Context, userID, sessionID string) error
}
