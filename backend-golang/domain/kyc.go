package domain

import (
	"context"
	"time"
)

type KYCDocument struct {
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type KYCStatus struct {
	UserID          string        `json:"user_id"`
	Status          string        `json:"kyc_status"`
	SubmittedAt     *time.Time    `json:"submitted_at"`
	ReviewedBy      *string       `json:"reviewed_by"`
	ReviewedAt      *time.Time    `json:"reviewed_at"`
	Documents       []KYCDocument `json:"documents"`
	RejectionReason *string       `json:"rejection_reason"`
}

type KYCUsecase interface {
	SubmitKYC(ctx context.Context, userID string, documents []KYCDocument) (*KYCStatus, error)
	GetKYCStatus(ctx context.Context, userID string) (*KYCStatus, error)
	ReviewKYC(ctx context.Context, userID, operator, status, note string) (*KYCStatus, error)
	IssueUploadToken(ctx context.Context, userID string) (map[string]interface{}, error)
}
