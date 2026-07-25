package usecases

import (
	"context"
	"testing"

	"folksart-be/backend-golang/internal/domain"
)

func TestRoleUsecase_Operations(t *testing.T) {
	u := NewRoleUsecase()
	ctx := context.Background()

	roles, err := u.ListRoles(ctx)
	if err != nil || len(roles) == 0 {
		t.Errorf("expected non-empty role list, got err: %v", err)
	}

	created, err := u.CreateRole(ctx, domain.Role{Name: "Auditor", Description: "Read-only audit"})
	if err != nil || created == nil || created.ID != "role-new" {
		t.Errorf("expected created role with ID role-new, got %v", created)
	}

	perms, err := u.ListPermissions(ctx)
	if err != nil || len(perms) == 0 {
		t.Errorf("expected permissions list, got err: %v", err)
	}

	if err := u.AssignRole(ctx, "usr-1", "role-1", "admin"); err != nil {
		t.Errorf("expected nil error on AssignRole, got %v", err)
	}

	if err := u.RemoveRole(ctx, "usr-1", "role-1", "admin"); err != nil {
		t.Errorf("expected nil error on RemoveRole, got %v", err)
	}
}
