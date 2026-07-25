package usecases

import (
	"context"
	"testing"

	"folksart-be/backend-golang/internal/domain"
)

type mockUserRepoForAuth struct {
	domain.UserRepository
}

func (m *mockUserRepoForAuth) Store(ctx context.Context, u *domain.IAMUser) error {
	return nil
}

func TestAuthUsecase_Login(t *testing.T) {
	u := NewAuthUsecase(&mockUserRepoForAuth{})
	ctx := context.Background()

	t.Run("Empty credentials", func(t *testing.T) {
		res, err := u.Login(ctx, "", "")
		if err == nil {
			t.Errorf("expected error for empty credentials, got nil")
		}
		if res != nil {
			t.Errorf("expected nil token response, got %v", res)
		}
	})

	t.Run("Valid credentials", func(t *testing.T) {
		res, err := u.Login(ctx, "admin@example.com", "secret123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.AccessToken == "" {
			t.Errorf("expected valid token response, got %v", res)
		}
		if res.User.Email != "admin@example.com" {
			t.Errorf("expected email admin@example.com, got %s", res.User.Email)
		}
	})
}

func TestAuthUsecase_Me(t *testing.T) {
	u := NewAuthUsecase(&mockUserRepoForAuth{})
	ctx := context.Background()

	user, err := u.Me(ctx, "sample-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil || user.ID != "usr-admin" {
		t.Errorf("expected usr-admin, got %v", user)
	}
}

func TestAuthUsecase_Sessions(t *testing.T) {
	u := NewAuthUsecase(&mockUserRepoForAuth{})
	ctx := context.Background()

	sessions, err := u.GetSessions(ctx, "usr-admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sessions) == 0 {
		t.Errorf("expected at least 1 session, got 0")
	}

	err = u.TerminateSession(ctx, "usr-admin", "sess-1")
	if err != nil {
		t.Errorf("expected nil error on session termination, got %v", err)
	}
}

func TestAuthUsecase_LogoutAndRefresh(t *testing.T) {
	u := NewAuthUsecase(&mockUserRepoForAuth{})
	ctx := context.Background()

	if err := u.Logout(ctx, "token"); err != nil {
		t.Errorf("expected nil error on logout, got %v", err)
	}

	_, err := u.Refresh(ctx, "token")
	if err != nil {
		t.Errorf("expected nil error on refresh, got %v", err)
	}
}

func TestAuthUsecase_RegisterCustomer(t *testing.T) {
	u := NewAuthUsecase(&mockUserRepoForAuth{})
	ctx := context.Background()

	t.Run("Missing credentials", func(t *testing.T) {
		res, err := u.RegisterCustomer(ctx, "", "", "", "")
		if err == nil || res != nil {
			t.Errorf("expected error for empty customer register credentials, got nil")
		}
	})

	t.Run("Successful CIAM registration", func(t *testing.T) {
		res, err := u.RegisterCustomer(ctx, "Jane Consumer", "jane@example.com", "pass123", "+62812345678")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil || res.User == nil {
			t.Fatalf("expected valid user response on customer register")
		}
		if res.User.UserType != "customer" || res.User.Role != "VerifiedCustomer" {
			t.Errorf("expected customer type with VerifiedCustomer role, got %v", res.User)
		}
		if res.User.TenantID != "public-ciam" {
			t.Errorf("expected tenantID public-ciam, got %s", res.User.TenantID)
		}
	})
}
