package dto

type UserResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Status     string `json:"status"`
	KYCStatus  string `json:"kycStatus"`
	Department string `json:"department"`
	RiskScore  int    `json:"riskScore"`
	UserType   string `json:"userType"`
	TenantID   string `json:"tenantId"`
	IsVerified bool   `json:"isVerified"`
	CreatedAt  string `json:"createdAt"`
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
	KYCStatus  *string `json:"kycStatus"`
	RiskScore  *int    `json:"riskScore"`
	MFAEnabled *bool   `json:"mfaEnabled"`
	Operator   string  `json:"operator" validate:"required"`
}
