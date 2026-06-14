package domain

import (
	"context"
	"time"
)

type AccessRequest struct {
	ID             string     `json:"id"`
	RequesterID    string     `json:"requester_id"`
	RequesterName  string     `json:"requester_name"`
	Resource       string     `json:"resource"`
	AccessLevel    string     `json:"access_level"`
	Justification  string     `json:"justification"`
	Status         string     `json:"status"` // pending | approved | rejected
	RequestedAt    time.Time  `json:"requested_at"`
	ApprovedBy     *string    `json:"approved_by,omitempty"`
	ApprovedAt     *time.Time `json:"approved_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
}

type AccessRequestUsecase interface {
	SubmitRequest(ctx context.Context, req AccessRequest) (*AccessRequest, error)
	ListRequests(ctx context.Context, status string) ([]AccessRequest, error)
	ApproveRequest(ctx context.Context, id, operator, note string) (*AccessRequest, error)
	RejectRequest(ctx context.Context, id, operator, reason string) (*AccessRequest, error)
}
