package usecases

import (
	"context"
	"time"
	"react-example/backend-golang/internal/domain"
)

type policyUsecase struct{}

func NewPolicyUsecase() domain.PolicyUsecase {
	return &policyUsecase{}
}

func (u *policyUsecase) ListPolicies(ctx context.Context) ([]domain.Policy, error) {
	return []domain.Policy{}, nil
}

func (u *policyUsecase) CreatePolicy(ctx context.Context, p domain.Policy) (*domain.Policy, error) {
	p.ID = "pol-123"
	p.CreatedAt = time.Now()
	return &p, nil
}

func (u *policyUsecase) UpdatePolicy(ctx context.Context, id string, p domain.Policy) (*domain.Policy, error) {
	return &p, nil
}

func (u *policyUsecase) DeletePolicy(ctx context.Context, id string) error {
	return nil
}

func (u *policyUsecase) Evaluate(ctx context.Context, userID, resource, action string) (*domain.PolicyEvaluation, error) {
	return &domain.PolicyEvaluation{Decision: "allow"}, nil
}
