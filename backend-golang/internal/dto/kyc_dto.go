package dto

type KYCDocumentDTO struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	UploadedAt string `json:"uploaded_at"`
}

type KYCStatusResponse struct {
	UserID          string           `json:"user_id"`
	Status          string           `json:"kyc_status"`
	SubmittedAt     *string          `json:"submitted_at,omitempty"`
	ReviewedBy      *string          `json:"reviewed_by,omitempty"`
	ReviewedAt      *string          `json:"reviewed_at,omitempty"`
	Documents       []KYCDocumentDTO `json:"documents"`
	RejectionReason *string          `json:"rejection_reason,omitempty"`
}

type KYCReviewRequest struct {
	Operator string `json:"operator" validate:"required"`
	Status   string `json:"status" validate:"required"`
	Note     string `json:"note"`
}
