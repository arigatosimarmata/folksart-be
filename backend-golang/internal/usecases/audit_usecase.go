package usecases

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"react-example/backend-golang/internal/domain"
)

type auditUsecase struct {
	auditRepo domain.AuditRepository
}

func NewAuditUsecase(ar domain.AuditRepository) domain.AuditUsecase {
	return &auditUsecase{auditRepo: ar}
}

func (u *auditUsecase) RecordAction(ctx context.Context, actor, action, target, severity string) (*domain.AuditLog, error) {
	if action == "" || target == "" || severity == "" {
		return nil, fmt.Errorf("insufficient log fields provided")
	}

	cleanActor := actor
	if cleanActor == "" {
		cleanActor = "anonymous_operator"
	}

	rand.Seed(time.Now().UnixNano())
	logID := fmt.Sprintf("log-%d", 1000+rand.Intn(9000))
	timestamp := time.Now()

	newLog := &domain.AuditLog{
		ID:        logID,
		Timestamp: timestamp,
		Actor:     cleanActor,
		Action:    action,
		Target:    target,
		Severity:  severity,
	}

	if err := u.auditRepo.Store(ctx, newLog); err != nil {
		return nil, err
	}

	return newLog, nil
}

func (u *auditUsecase) GetAuditTrail(ctx context.Context, severity string, limit int) ([]domain.AuditLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return u.auditRepo.List(ctx, severity, limit)
}
