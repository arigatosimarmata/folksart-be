package domain

import (
	"context"
	"time"
)

type NotificationRule struct {
	ID        string
	Name      string
	Trigger   string
	Severity  string
	Channels  []string
	Active    bool
	CreatedAt time.Time
}

type Notification struct {
	ID        string
	Type      string
	Title     string
	Body      string
	Severity  string
	Read      bool
	CreatedAt time.Time
}

type NotificationUsecase interface {
	ListRules(ctx context.Context) ([]NotificationRule, error)
	CreateRule(ctx context.Context, rule NotificationRule) (*NotificationRule, error)
	ListNotifications(ctx context.Context, userID string) ([]Notification, error)
}
