package usecase

import (
	"context"
	"time"
	"react-example/backend-golang/domain"
)

type policyUsecase struct{}

func NewPolicyUsecase() domain.PolicyUsecase {
	return &policyUsecase{}
}

func (u *policyUsecase) ListPolicies(ctx context.Context) ([]domain.Policy, error) {
	return []domain.Policy{
		{
			ID:          "pol_01HRISK",
			Name:        "high-risk-block",
			Description: "Blokir akses jika risk_score > 80",
			Condition:   map[string]interface{}{"field": "risk_score", "operator": "gt", "value": 80},
			Action:      "deny",
			Priority:    1,
			Active:      true,
			CreatedAt:   time.Now().AddDate(0, -1, 0),
		},
		{
			ID:          "pol_02HKYC",
			Name:        "unverified-kyc-restrict",
			Description: "Batasi akses jika KYC belum verified",
			Condition:   map[string]interface{}{"field": "kyc_status", "operator": "neq", "value": "verified"},
			Action:      "restrict",
			Priority:    2,
			Active:      true,
			CreatedAt:   time.Now().AddDate(0, -1, 0),
		},
	}, nil
}

func (u *policyUsecase) CreatePolicy(ctx context.Context, p domain.Policy) (*domain.Policy, error) {
	p.ID = "pol_new_" + time.Now().Format("05")
	p.CreatedAt = time.Now()
	return &p, nil
}

func (u *policyUsecase) UpdatePolicy(ctx context.Context, id string, p domain.Policy) (*domain.Policy, error) {
	p.ID = id
	return &p, nil
}

func (u *policyUsecase) DeletePolicy(ctx context.Context, id string) error {
	return nil
}

func (u *policyUsecase) Evaluate(ctx context.Context, userID, resource, action string) (*domain.PolicyEvaluation, error) {
	// Mock decision: deny if risk score > 80
	return &domain.PolicyEvaluation{
		UserID:   userID,
		Resource: resource,
		Action:   action,
		Decision: "allow",
	}, nil
}
