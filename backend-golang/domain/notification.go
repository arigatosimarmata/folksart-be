package domain

import (
	"context"
	"time"
)

type NotificationRule struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Trigger   string   `json:"trigger"`
	Severity  string   `json:"severity"`
	Channels  []string `json:"channels"`
	Active    bool     `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Severity  string    `json:"severity"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationUsecase interface {
	ListRules(ctx context.Context) ([]NotificationRule, error)
	CreateRule(ctx context.Context, rule NotificationRule) (*NotificationRule, error)
	ListNotifications(ctx context.Context, userID string) ([]Notification, error)
}
