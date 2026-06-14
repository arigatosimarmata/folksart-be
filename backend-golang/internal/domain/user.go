package domain

import (
	"context"
	"time"
)

type IAMUser struct {
	ID          string
	Name        string
	Username    string
	Email       string
	Phone       string
	Role        string
	Status      string
	KYCStatus   string
	Department  string
	RiskScore   int
	MFAEnabled  bool
	CreatedAt   time.Time
}

type UserFilter struct {
	Search     string
	Role       string
	Status     string
	Department string
	Limit      int
	Offset     int
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*IAMUser, error)
	List(ctx context.Context, filter UserFilter) ([]IAMUser, int, error)
	Store(ctx context.Context, user *IAMUser) error
	Update(ctx context.Context, user *IAMUser) error
	Delete(ctx context.Context, id string) error
}

type UserUsecase interface {
	FetchDirectories(ctx context.Context, filter UserFilter) ([]IAMUser, int, error)
	EnrollPrincipal(ctx context.Context, name, email, role, department, operator string) (*IAMUser, error)
	PatchPrincipal(ctx context.Context, id string, status, kycStatus *string, riskScore *int, mfaEnabled *bool, operator string) (*IAMUser, error)
	DecommissionPrincipal(ctx context.Context, id, operator string) error
	ExportCSVStream(ctx context.Context, filter UserFilter) ([]IAMUser, error)
}
