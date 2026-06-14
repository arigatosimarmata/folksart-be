package usecase

import (
	"context"
	"time"
	"react-example/backend-golang/domain"
)

type accessRequestUsecase struct{}

func NewAccessRequestUsecase() domain.AccessRequestUsecase {
	return &accessRequestUsecase{}
}

func (u *accessRequestUsecase) SubmitRequest(ctx context.Context, req domain.AccessRequest) (*domain.AccessRequest, error) {
	req.ID = "req_01HREQ5678"
	req.Status = "pending"
	req.RequestedAt = time.Now()
	return &req, nil
}

func (u *accessRequestUsecase) ListRequests(ctx context.Context, status string) ([]domain.AccessRequest, error) {
	return []domain.AccessRequest{
		{
			ID:            "req_01HREQ5678",
			RequesterName:  "Jane Smith",
			Resource:       "finance-reports",
			AccessLevel:    "read",
			Status:         "pending",
			RequestedAt:    time.Now().Add(-1 * time.Hour),
		},
		{
			ID:            "req_00HREQ1111",
			RequesterName:  "Budi Santoso",
			Resource:       "audit-logs",
			AccessLevel:    "read",
			Status:         "approved",
			RequestedAt:    time.Now().Add(-48 * time.Hour),
		},
	}, nil
}

func (u *accessRequestUsecase) ApproveRequest(ctx context.Context, id, operator, note string) (*domain.AccessRequest, error) {
	exp := time.Now().AddDate(0, 1, 0)
	now := time.Now()
	return &domain.AccessRequest{
		ID:         id,
		Status:     "approved",
		ApprovedBy: &operator,
		ApprovedAt: &now,
		ExpiresAt:  &exp,
	}, nil
}

func (u *accessRequestUsecase) RejectRequest(ctx context.Context, id, operator, reason string) (*domain.AccessRequest, error) {
	now := time.Now()
	return &domain.AccessRequest{
		ID:         id,
		Status:     "rejected",
		ApprovedBy: &operator,
		ApprovedAt: &now,
	}, nil
}
