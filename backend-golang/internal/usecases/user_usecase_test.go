package usecases

import (
	"context"
	"errors"
	"testing"

	"folksart-be/backend-golang/internal/domain"
)

type mockUserRepo struct {
	users  map[string]*domain.IAMUser
	errGet error
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.IAMUser, error) {
	if m.errGet != nil {
		return nil, m.errGet
	}
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	// return copy
	u := *user
	return &u, nil
}

func (m *mockUserRepo) List(ctx context.Context, filter domain.UserFilter) ([]domain.IAMUser, int, error) {
	var res []domain.IAMUser
	for _, u := range m.users {
		res = append(res, *u)
	}
	return res, len(res), nil
}

func (m *mockUserRepo) Store(ctx context.Context, user *domain.IAMUser) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(ctx context.Context, user *domain.IAMUser) error {
	if _, ok := m.users[user.ID]; !ok {
		return errors.New("user not found in update")
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id string) error {
	delete(m.users, id)
	return nil
}

// mockAuditRepo verifies User Global Rules: every UPDATE and DELETE MUST record an Alert/Audit notification
type mockAuditRepo struct {
	storedLogs []*domain.AuditLog
}

func (m *mockAuditRepo) Store(ctx context.Context, log *domain.AuditLog) error {
	m.storedLogs = append(m.storedLogs, log)
	return nil
}

func (m *mockAuditRepo) List(ctx context.Context, severity string, limit int) ([]domain.AuditLog, error) {
	return nil, nil
}

func TestUserUsecase_EnrollAndDirectories(t *testing.T) {
	userRepo := &mockUserRepo{users: make(map[string]*domain.IAMUser)}
	auditRepo := &mockAuditRepo{}
	u := NewUserUsecase(userRepo, auditRepo)
	ctx := context.Background()

	user, err := u.EnrollPrincipal(ctx, "John Doe", "john@example.com", "Analyst", "Finance", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil || user.Name != "John Doe" {
		t.Errorf("expected John Doe created, got %v", user)
	}

	// Ensure audit trail is recorded on store
	if len(auditRepo.storedLogs) != 1 {
		t.Errorf("expected 1 audit log created during enrollment, got %d", len(auditRepo.storedLogs))
	}

	dirs, count, err := u.FetchDirectories(ctx, domain.UserFilter{})
	if err != nil || count != 1 || len(dirs) != 1 {
		t.Errorf("expected 1 user in directory, got %d", count)
	}
}

func TestUserUsecase_PatchPrincipal_AuditAlert(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]*domain.IAMUser{
			"usr-101": {ID: "usr-101", Name: "Alice", Status: "Active", RiskScore: 20},
		},
	}
	auditRepo := &mockAuditRepo{}
	u := NewUserUsecase(userRepo, auditRepo)
	ctx := context.Background()

	newStatus := "Banned"
	newScore := 90
	updated, err := u.PatchPrincipal(ctx, "usr-101", &newStatus, nil, &newScore, nil, "security-officer")
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if updated.Status != "Banned" || updated.RiskScore != 90 {
		t.Errorf("expected status Banned and risk 90, got %v", updated)
	}

	// Verify Aturan Global (Global Rules): Update database MUST have notification/alert query log!
	if len(auditRepo.storedLogs) != 1 {
		t.Fatalf("CRITICAL RULE VIOLATION: expected alert/audit notification on update, got 0 logs")
	}
	log := auditRepo.storedLogs[0]
	if log.Severity != "Critical" {
		t.Errorf("expected Critical severity for banning account, got %s", log.Severity)
	}
}

func TestUserUsecase_DecommissionPrincipal_AuditAlert(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]*domain.IAMUser{
			"usr-202": {ID: "usr-202", Name: "Bob", Status: "Active"},
		},
	}
	auditRepo := &mockAuditRepo{}
	u := NewUserUsecase(userRepo, auditRepo)
	ctx := context.Background()

	err := u.DecommissionPrincipal(ctx, "usr-202", "admin-1")
	if err != nil {
		t.Fatalf("unexpected error during decommissioning: %v", err)
	}

	// Verify Aturan Global (Global Rules): Delete database MUST generate alert log!
	if len(auditRepo.storedLogs) != 1 {
		t.Fatalf("CRITICAL RULE VIOLATION: expected alert/audit notification on delete, got 0 logs")
	}
	log := auditRepo.storedLogs[0]
	if log.Severity != "Critical" || log.Action != "Permanent Governance Offboarding (Account Decommissioned)" {
		t.Errorf("expected Critical decommission log, got %v", log)
	}
}

func TestUserUsecase_ExportCSVStream(t *testing.T) {
	userRepo := &mockUserRepo{
		users: map[string]*domain.IAMUser{
			"usr-1": {ID: "usr-1", Name: "User 1"},
		},
	}
	u := NewUserUsecase(userRepo, &mockAuditRepo{})
	ctx := context.Background()

	stream, err := u.ExportCSVStream(ctx, domain.UserFilter{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stream) != 1 {
		t.Errorf("expected 1 user exported in CSV stream, got %d", len(stream))
	}
}
