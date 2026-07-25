package dto

type KYCDocumentDTO struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	UploadedAt string `json:"uploadedAt"`
}

type KYCStatusResponse struct {
	UserID          string           `json:"userId"`
	Status          string           `json:"kycStatus"`
	SubmittedAt     *string          `json:"submittedAt,omitempty"`
	ReviewedBy      *string          `json:"reviewedBy,omitempty"`
	ReviewedAt      *string          `json:"reviewedAt,omitempty"`
	Documents       []KYCDocumentDTO `json:"documents"`
	RejectionReason *string          `json:"rejectionReason,omitempty"`
}

type KYCReviewRequest struct {
	Operator string `json:"operator" validate:"required"`
	Status   string `json:"status" validate:"required"`
	Note     string `json:"note"`
}
