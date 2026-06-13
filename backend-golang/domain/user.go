package domain

import (
	"context"
	"time"
)

// IAMUser defines the core domain entity representing corporate identities
type IAMUser struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Role        string    `json:"role"`        // Administrator | Security Officer | End User
	Status      string    `json:"status"`      // Active | Banned | Deactivated
	KYCStatus   string    `json:"kycStatus"`   // Pending | Verified | Failed | Suspicious
	Department  string    `json:"department"`
	RiskScore   int       `json:"riskScore"`   // 0 to 100
	MFAEnabled  bool      `json:"mfaEnabled"`
	CreatedAt   time.Time `json:"createdAt"`
}

// UserFilter defines the schema for query matrices
type UserFilter struct {
	Search     string
	Role       string
	Status     string
	Department string
	Limit      int
	Offset     int
}

// UserRepository specifies fine-grained repository interfaces supporting Interface Segregation
type UserRepository interface {
	GetByID(ctx context.Context, id string) (*IAMUser, error)
	List(ctx context.Context, filter UserFilter) ([]IAMUser, int, error)
	Store(ctx context.Context, user *IAMUser) error
	Update(ctx context.Context, user *IAMUser) error
	Delete(ctx context.Context, id string) error
}

// UserUsecase defines application boundary contracts
type UserUsecase interface {
	FetchDirectories(ctx context.Context, filter UserFilter) ([]IAMUser, int, error)
	EnrollPrincipal(ctx context.Context, name, email, role, department, operator string) (*IAMUser, error)
	PatchPrincipal(ctx context.Context, id string, status, kycStatus *string, riskScore *int, mfaEnabled *bool, operator string) (*IAMUser, error)
	DecommissionPrincipal(ctx context.Context, id, operator string) error
	ExportCSVStream(ctx context.Context, filter UserFilter) ([]IAMUser, error)
}
