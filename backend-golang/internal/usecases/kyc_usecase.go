package usecases

import (
	"context"
	"time"
	"react-example/backend-golang/internal/domain"
)

type kycUsecase struct{}

func NewKYCUsecase() domain.KYCUsecase {
	return &kycUsecase{}
}

func (u *kycUsecase) SubmitKYC(ctx context.Context, userID string, documents []domain.KYCDocument) (*domain.KYCStatus, error) {
	now := time.Now()
	return &domain.KYCStatus{
		UserID:      userID,
		Status:      "Pending",
		SubmittedAt: &now,
		Documents:   documents,
	}, nil
}

func (u *kycUsecase) GetKYCStatus(ctx context.Context, userID string) (*domain.KYCStatus, error) {
	return &domain.KYCStatus{UserID: userID, Status: "Verified"}, nil
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
	return map[string]interface{}{"upload_token": "token-123"}, nil
}
