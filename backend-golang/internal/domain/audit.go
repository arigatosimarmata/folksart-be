package domain

import (
	"context"
	"time"
)

type AuditLog struct {
	ID        string
	Timestamp time.Time
	Actor     string
	Action    string
	Target    string
	Severity  string
}

type AuditRepository interface {
	Store(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, severity string, limit int) ([]AuditLog, error)
}

type AuditUsecase interface {
	RecordAction(ctx context.Context, actor, action, target, severity string) (*AuditLog, error)
	GetAuditTrail(ctx context.Context, severity string, limit int) ([]AuditLog, error)
}
