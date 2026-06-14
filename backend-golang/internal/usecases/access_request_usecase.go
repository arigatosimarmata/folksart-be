package usecases

import (
	"context"
	"time"
	"react-example/backend-golang/internal/domain"
)

type arUsecase struct{}

func NewAccessRequestUsecase() domain.AccessRequestUsecase {
	return &arUsecase{}
}

func (u *arUsecase) SubmitRequest(ctx context.Context, req domain.AccessRequest) (*domain.AccessRequest, error) {
	req.ID = "req-123"
	req.Status = "pending"
	req.RequestedAt = time.Now()
	return &req, nil
}

func (u *arUsecase) ListRequests(ctx context.Context, status string) ([]domain.AccessRequest, error) {
	return []domain.AccessRequest{}, nil
}

func (u *arUsecase) ApproveRequest(ctx context.Context, id, operator, note string) (*domain.AccessRequest, error) {
	now := time.Now()
	return &domain.AccessRequest{
		ID:         id,
		Status:     "approved",
		ApprovedBy: &operator,
		ApprovedAt: &now,
	}, nil
}

func (u *arUsecase) RejectRequest(ctx context.Context, id, operator, reason string) (*domain.AccessRequest, error) {
	return &domain.AccessRequest{
		ID:     id,
		Status: "rejected",
	}, nil
}
