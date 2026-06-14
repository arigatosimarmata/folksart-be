package usecase

import (
	"context"
	"time"
	"react-example/backend-golang/domain"
)

type notificationUsecase struct{}

func NewNotificationUsecase() domain.NotificationUsecase {
	return &notificationUsecase{}
}

func (u *notificationUsecase) ListRules(ctx context.Context) ([]domain.NotificationRule, error) {
	return []domain.NotificationRule{
		{
			ID:        "rule_01H",
			Name:      "new-ip-login-alert",
			Trigger:   "LOGIN_NEW_IP",
			Severity:  "Warning",
			Channels:  []string{"in_app", "email"},
			Active:    true,
			CreatedAt: time.Now().AddDate(0, -6, 0),
		},
	}, nil
}

func (u *notificationUsecase) CreateRule(ctx context.Context, rule domain.NotificationRule) (*domain.NotificationRule, error) {
	rule.ID = "rule_new"
	rule.CreatedAt = time.Now()
	return &rule, nil
}

func (u *notificationUsecase) ListNotifications(ctx context.Context, userID string) ([]domain.Notification, error) {
	return []domain.Notification{
		{
			ID:        "notif_01H",
			Type:      "SECURITY_ALERT",
			Title:     "Login dari IP baru terdeteksi",
			Body:      "Login berhasil dari 203.0.113.99 pada 10 Jun 02:15.",
			Severity:  "Warning",
			Read:      false,
			CreatedAt: time.Now().Add(-12 * time.Hour),
		},
		{
			ID:        "notif_02H",
			Type:      "ACCESS_REQUEST",
			Title:     "Request akses membutuhkan persetujuan",
			Body:      "Jane Smith mengajukan akses ke finance-reports.",
			Severity:  "Info",
			Read:      true,
			CreatedAt: time.Now().Add(-5 * time.Hour),
		},
	}, nil
}
