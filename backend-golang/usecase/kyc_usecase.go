package usecase

import (
	"context"
	"time"
	"react-example/backend-golang/domain"
)

type kycUsecase struct{}

func NewKYCUsecase() domain.KYCUsecase {
	return &kycUsecase{}
}

func (u *kycUsecase) SubmitKYC(ctx context.Context, userID string, docs []domain.KYCDocument) (*domain.KYCStatus, error) {
	now := time.Now()
	for i := range docs {
		docs[i].Status = "uploaded"
		docs[i].UploadedAt = now
	}
	return &domain.KYCStatus{
		UserID:      userID,
		Status:      "under_review",
		SubmittedAt: &now,
		Documents:   docs,
	}, nil
}

func (u *kycUsecase) GetKYCStatus(ctx context.Context, userID string) (*domain.KYCStatus, error) {
	now := time.Now().AddDate(0, 0, -1)
	return &domain.KYCStatus{
		UserID:      userID,
		Status:      "under_review",
		SubmittedAt: &now,
		Documents: []domain.KYCDocument{
			{Type: "id_card", Status: "uploaded", UploadedAt: now},
			{Type: "selfie", Status: "uploaded", UploadedAt: now},
		},
	}, nil
}

func (u *kycUsecase) ReviewKYC(ctx context.Context, userID, operator, status, note string) (*domain.KYCStatus, error) {
	now := time.Now()
	return &domain.KYCStatus{
		UserID:     userID,
		Status:     status,
		ReviewedBy: &operator,
		ReviewedAt: &now,
	}, nil
}

func (u *kycUsecase) IssueUploadToken(ctx context.Context, userID string) (map[string]interface{}, error) {
	return map[string]interface{}{
		"upload_token":  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		"token_type":    "JWT",
		"expires_in":    600,
		"upload_url":    "https://storage.example.com/kyc-docs/" + userID + "/",
		"allowed_types": []string{"image/jpeg", "image/png", "application/pdf"},
		"max_size_mb":   10,
	}, nil
}
