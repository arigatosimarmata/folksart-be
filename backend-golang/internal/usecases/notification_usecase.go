package usecases

import (
	"context"
	"react-example/backend-golang/internal/domain"
)

type nfUsecase struct{}

func NewNotificationUsecase() domain.NotificationUsecase {
	return &nfUsecase{}
}

func (u *nfUsecase) ListRules(ctx context.Context) ([]domain.NotificationRule, error) {
	return []domain.NotificationRule{}, nil
}

func (u *nfUsecase) CreateRule(ctx context.Context, rule domain.NotificationRule) (*domain.NotificationRule, error) {
	rule.ID = "rule-123"
	return &rule, nil
}

func (u *nfUsecase) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	return []domain.Notification{}, nil
}
