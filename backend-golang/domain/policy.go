package domain

import (
	"context"
	"time"
)

type Policy struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Condition   interface{} `json:"condition"`
	Action      string      `json:"action"` // allow | deny | restrict
	Priority    int         `json:"priority"`
	Active      bool        `json:"active"`
	CreatedAt   time.Time   `json:"created_at"`
}

type PolicyEvaluation struct {
	UserID        string `json:"user_id"`
	Resource      string `json:"resource"`
	Action        string `json:"action"`
	Decision      string `json:"decision"`
	MatchedPolicy *Policy `json:"matched_policy,omitempty"`
}

type PolicyUsecase interface {
	ListPolicies(ctx context.Context) ([]Policy, error)
	CreatePolicy(ctx context.Context, p Policy) (*Policy, error)
	UpdatePolicy(ctx context.Context, id string, p Policy) (*Policy, error)
	DeletePolicy(ctx context.Context, id string) error
	Evaluate(ctx context.Context, userID, resource, action string) (*PolicyEvaluation, error)
}
