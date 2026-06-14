package domain

import (
	"context"
	"time"
)

type AccessRequest struct {
	ID             string
	RequesterID    string
	RequesterName  string
	Resource       string
	AccessLevel    string
	Justification  string
	Status         string
	RequestedAt    time.Time
	ApprovedBy     *string
	ApprovedAt     *time.Time
	ExpiresAt      *time.Time
}

type AccessRequestUsecase interface {
	SubmitRequest(ctx context.Context, req AccessRequest) (*AccessRequest, error)
	ListRequests(ctx context.Context, status string) ([]AccessRequest, error)
	ApproveRequest(ctx context.Context, id, operator, note string) (*AccessRequest, error)
	RejectRequest(ctx context.Context, id, operator, reason string) (*AccessRequest, error)
}
