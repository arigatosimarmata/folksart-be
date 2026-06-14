package dto

type UserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	KYCStatus  string `json:"kyc_status"`
	Department string `json:"department"`
	RiskScore  int    `json:"risk_score"`
	CreatedAt  string `json:"created_at"`
}

type EnrollUserRequest struct {
	Name       string `json:"name" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Role       string `json:"role" validate:"required"`
	Department string `json:"department" validate:"required"`
	Operator   string `json:"operator" validate:"required"`
}

type UpdateUserRequest struct {
	Status     *string `json:"status"`
	KYCStatus  *string `json:"kyc_status"`
	RiskScore  *int    `json:"risk_score"`
	MFAEnabled *bool   `json:"mfa_enabled"`
	Operator   string  `json:"operator" validate:"required"`
}
