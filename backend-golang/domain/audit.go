package domain

import (
	"context"
	"time"
)

// AuditLog acts as the records representing immutable administrative action trails
type AuditLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Actor     string    `json:"actor"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Severity  string    `json:"severity"` // Critical | High | Medium | Low
}

// AuditRepository contract defining write and read schemas for audit trails
type AuditRepository interface {
	Store(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, severity string, limit int) ([]AuditLog, error)
}

// AuditUsecase contract defining validation and entry orchestration rules
type AuditUsecase interface {
	RecordAction(ctx context.Context, actor, action, target, severity string) (*AuditLog, error)
	GetAuditTrail(ctx context.Context, severity string, limit int) ([]AuditLog, error)
}
