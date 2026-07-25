package usecases

import (
	"context"
	"testing"

	"folksart-be/backend-golang/internal/domain"
)

func TestPolicyUsecase_Operations(t *testing.T) {
	u := NewPolicyUsecase()
	ctx := context.Background()

	policies, err := u.ListPolicies(ctx)
	if err != nil || policies == nil {
		t.Errorf("expected policies slice, got err: %v", err)
	}

	newPol, err := u.CreatePolicy(ctx, domain.Policy{Name: "No-Access-On-Weekends"})
	if err != nil || newPol.ID != "pol-123" {
		t.Errorf("expected policy created with id pol-123, got %v", newPol)
	}

	updated, err := u.UpdatePolicy(ctx, "pol-123", domain.Policy{Name: "Updated-Policy"})
	if err != nil || updated.Name != "Updated-Policy" {
		t.Errorf("expected policy updated, got %v", updated)
	}

	if err := u.DeletePolicy(ctx, "pol-123"); err != nil {
		t.Errorf("expected nil error on delete, got %v", err)
	}

	eval, err := u.Evaluate(ctx, "usr-1", "report", "read")
	if err != nil || eval.Decision != "allow" {
		t.Errorf("expected evaluation decision allow, got %v", eval)
	}
}
