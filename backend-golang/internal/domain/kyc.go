package domain

import (
	"context"
	"time"
)

type KYCDocument struct {
	Type       string
	Status     string
	UploadedAt time.Time
}

type KYCStatus struct {
	UserID          string
	Status          string
	SubmittedAt     *time.Time
	ReviewedBy      *string
	ReviewedAt      *time.Time
	Documents       []KYCDocument
	RejectionReason *string
}

type KYCUsecase interface {
	SubmitKYC(ctx context.Context, userID string, documents []KYCDocument) (*KYCStatus, error)
	GetKYCStatus(ctx context.Context, userID string) (*KYCStatus, error)
	ReviewKYC(ctx context.Context, userID, operator, status, note string) (*KYCStatus, error)
	IssueUploadToken(ctx context.Context, userID string) (map[string]interface{}, error)
}
