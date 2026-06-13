package usecase

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"react-example/backend-golang/domain"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	auditRepo domain.AuditRepository
}

// NewUserUsecase initializes the application interactor for Users
func NewUserUsecase(ur domain.UserRepository, ar domain.AuditRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo:  ur,
		auditRepo: ar,
	}
}

func (u *userUsecase) FetchDirectories(ctx context.Context, filter domain.UserFilter) ([]domain.IAMUser, int, error) {
	return u.userRepo.List(ctx, filter)
}

func (u *userUsecase) EnrollPrincipal(ctx context.Context, name, email, role, department, operator string) (*domain.IAMUser, error) {
	if name == "" || email == "" || role == "" || department == "" {
		return nil, fmt.Errorf("insufficient enrollment credentials provided")
	}

	cleanUsername := strings.ToLower(strings.ReplaceAll(name, " ", ""))
	if len(cleanUsername) > 15 {
		cleanUsername = cleanUsername[:15]
	}

	rand.Seed(time.Now().UnixNano())
	randomID := fmt.Sprintf("usr-%d", 1000+rand.Intn(9000))
	createdAt := time.Now()
	riskScore := 10 + rand.Intn(25) // Baseline risk score

	newUser := &domain.IAMUser{
		ID:         randomID,
		Name:       name,
		Username:   cleanUsername,
		Email:      email,
		Phone:      "+1 (555) 012-7492",
		Role:       role,
		Status:     "Active",
		KYCStatus:  "Pending",
		Department: department,
		RiskScore:  riskScore,
		MFAEnabled: false,
		CreatedAt:  createdAt,
	}

	// 1. Persist User Entity
	if err := u.userRepo.Store(ctx, newUser); err != nil {
		return nil, fmt.Errorf("usecase failed creating client record: %w", err)
	}

	// 2. Log Action entry
	logID := fmt.Sprintf("log-%d", 1000+rand.Intn(9000))
	actor := operator
	if actor == "" {
		actor = "identity_provisioner"
	}
	actionStr := fmt.Sprintf("Enrolled Subject Principal (Role: %s)", role)

	newAudit := &domain.AuditLog{
		ID:        logID,
		Timestamp: createdAt,
		Actor:     actor,
		Action:    actionStr,
		Target:    name,
		Severity:  "Medium",
	}

	_ = u.auditRepo.Store(ctx, newAudit) // Save but don't halt flow on async logging trace errors

	return newUser, nil
}

func (u *userUsecase) PatchPrincipal(ctx context.Context, id string, status, kycStatus *string, riskScore *int, mfaEnabled *bool, operator string) (*domain.IAMUser, error) {
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed looking up target user: %w", err)
	}
	if existingUser == nil {
		return nil, fmt.Errorf("user metadata not discovered matching identifier")
	}

	var changes []string

	if status != nil {
		changes = append(changes, fmt.Sprintf("Status altered from '%s' to '%s'", existingUser.Status, *status))
		existingUser.Status = *status
	}
	if kycStatus != nil {
		changes = append(changes, fmt.Sprintf("KYC updated to '%s'", *kycStatus))
		existingUser.KYCStatus = *kycStatus
	}
	if riskScore != nil {
		safeScore := *riskScore
		if safeScore < 0 {
			safeScore = 0
		} else if safeScore > 100 {
			safeScore = 100
		}
		changes = append(changes, fmt.Sprintf("Risk shifted from %d%% to %d%%", existingUser.RiskScore, safeScore))
		existingUser.RiskScore = safeScore
	}
	if mfaEnabled != nil {
		changes = append(changes, fmt.Sprintf("MFA switch shifted to %v", *mfaEnabled))
		existingUser.MFAEnabled = *mfaEnabled
	}

	if len(changes) == 0 {
		return existingUser, nil // Nothing to update
	}

	// 1. Commit update
	if err := u.userRepo.Update(ctx, existingUser); err != nil {
		return nil, fmt.Errorf("usecase failed to update db models: %w", err)
	}

	// 2. Logging audit trace
	severity := "Low"
	if status != nil && *status == "Banned" {
		severity = "Critical"
	} else if riskScore != nil && *riskScore > 75 {
		severity = "High"
	} else if kycStatus != nil {
		severity = "Medium"
	}

	logID := fmt.Sprintf("log-%d", 1000+rand.Intn(9000))
	actor := operator
	if actor == "" {
		actor = "security_officer"
	}
	actionText := fmt.Sprintf("Attributes Patch: [%s]", strings.Join(changes, " | "))

	audit := &domain.AuditLog{
		ID:        logID,
		Timestamp: time.Now(),
		Actor:     actor,
		Action:    actionText,
		Target:    existingUser.Name,
		Severity:  severity,
	}

	_ = u.auditRepo.Store(ctx, audit)

	return existingUser, nil
}

func (u *userUsecase) DecommissionPrincipal(ctx context.Context, id, operator string) error {
	existingUser, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("decommission query lookup failed: %w", err)
	}
	if existingUser == nil {
		return fmt.Errorf("identity record not found")
	}

	// 1. Purge identity
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("decommission execution failed: %w", err)
	}

	// 2. Write Audit statement
	logID := fmt.Sprintf("log-%d", 1000+rand.Intn(9000))
	actor := operator
	if actor == "" {
		actor = "admin_operator"
	}

	audit := &domain.AuditLog{
		ID:        logID,
		Timestamp: time.Now(),
		Actor:     actor,
		Action:    "Permanent Governance Offboarding (Account Decommissioned)",
		Target:    existingUser.Name,
		Severity:  "Critical",
	}

	_ = u.auditRepo.Store(ctx, audit)

	return nil
}

func (u *userUsecase) ExportCSVStream(ctx context.Context, filter domain.UserFilter) ([]domain.IAMUser, error) {
	// For csv compilation, we retrieve all users without limits (offset offset reset to zero/empty)
	csvFilter := filter
	csvFilter.Limit = 0
	csvFilter.Offset = 0

	users, _, err := u.userRepo.List(ctx, csvFilter)
	return users, err
}
