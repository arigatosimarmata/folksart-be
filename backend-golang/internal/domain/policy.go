package domain

import (
	"context"
	"time"
)

type Policy struct {
	ID          string
	Name        string
	Description string
	Condition   interface{}
	Action      string
	Priority    int
	Active      bool
	CreatedAt   time.Time
}

type PolicyEvaluation struct {
	UserID        string
	Resource      string
	Action        string
	Decision      string
	MatchedPolicy *Policy
}

type PolicyUsecase interface {
	ListPolicies(ctx context.Context) ([]Policy, error)
	CreatePolicy(ctx context.Context, p Policy) (*Policy, error)
	UpdatePolicy(ctx context.Context, id string, p Policy) (*Policy, error)
	DeletePolicy(ctx context.Context, id string) error
	Evaluate(ctx context.Context, userID, resource, action string) (*PolicyEvaluation, error)
}
