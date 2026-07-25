package security

import (
	"testing"
	"folksart-be/backend-golang/config"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "SecretP@ssw0rd!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if !CheckPasswordHash(password, hash) {
		t.Errorf("expected password check to succeed")
	}

	if CheckPasswordHash("WrongPassword", hash) {
		t.Errorf("expected wrong password check to fail")
	}
}

func TestJWTGenerateAndValidate(t *testing.T) {
	// Set custom secret
	config.AppConfig.JWTSecretKey = "super-secret-key-for-testing"

	tokenStr, err := GenerateToken("usr-admin", "Administrator")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}
	if tokenStr == "" {
		t.Fatalf("expected non-empty token string")
	}

	claims, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}
	if claims.UserID != "usr-admin" {
		t.Errorf("expected UserID usr-admin, got %s", claims.UserID)
	}
	if claims.Role != "Administrator" {
		t.Errorf("expected Role Administrator, got %s", claims.Role)
	}

	_, err = ValidateToken("invalid.token.string")
	if err == nil {
		t.Errorf("expected error on invalid token, got nil")
	}
}
